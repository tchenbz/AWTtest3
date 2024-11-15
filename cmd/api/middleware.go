package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
    "github.com/dgrijalva/jwt-go"
)

func (a *applicationDependencies)recoverPanic(next http.Handler)http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func ()  {
			err := recover();
			if err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *applicationDependencies) rateLimit(next http.Handler) http.Handler {
    type client struct {
        limiter  *rate.Limiter
        lastSeen time.Time
    }
    var mu sync.Mutex    
    var clients = make(map[string]*client)
    go func() {
        for {
            time.Sleep(time.Minute)
            mu.Lock() 
            for ip, client := range clients {
                if time.Since(client.lastSeen) > 3*time.Minute {
                    delete(clients, ip)
                }
            }
            mu.Unlock() 
        }
    }()
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if a.config.limiter.enabled {
            ip, _, err := net.SplitHostPort(r.RemoteAddr)
            if err != nil {
                a.serverErrorResponse(w, r, err)
                return
            }

            mu.Lock() 
            _, found := clients[ip]
            if !found {
                clients[ip] = &client{limiter: rate.NewLimiter(
                    rate.Limit(a.config.limiter.rps),
                    a.config.limiter.burst),
                }
            }
            clients[ip].lastSeen = time.Now()

            if !clients[ip].limiter.Allow() {
                mu.Unlock() 
                a.rateLimitExceededResponse(w, r)
                return
            }

            mu.Unlock()
        } 
        next.ServeHTTP(w, r)
    })

}

func CheckAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Proceed to next handler if token is valid
        next.ServeHTTP(w, r)
    })
}