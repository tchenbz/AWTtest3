package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/tchenbz/comments/internal/data"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port int
	environment string
	db struct {
        dsn string
    }
}

type applicationDependencies struct {
	config serverConfig
	logger *slog.Logger
	commentModel data.CommentModel
}

func main() {
	var settings serverConfig

	flag.IntVar(&settings.port, "port", 4000, "Server port")
	flag.StringVar(&settings.environment, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://quiz3:quiz3@localhost/comments?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// the call to openDB() sets up our connection pool
	db, err := openDB(settings)
	if err != nil {
    logger.Error(err.Error())
    os.Exit(1)
	}

	// release the database resources before exiting
	defer db.Close()

	logger.Info("database connection pool established")


	appInstance := &applicationDependencies {
		config: settings,
		logger: logger,
		commentModel: data.CommentModel{DB: db},
	}

	//router := http.NewServeMux()
	//router.HandleFunc("/v1/healthcheck", appInstance.healthCheckHandler)

	apiServer := &http.Server{ 
		Addr: fmt.Sprintf(":%d", settings.port),
		Handler: appInstance.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "address", apiServer.Addr, "environment", settings.environment)
	err = apiServer.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(settings serverConfig) (*sql.DB, error) {
    // open a connection pool
    db, err := sql.Open("postgres", settings.db.dsn)
    if err != nil {
        return nil, err
    }
    
    // set a context to ensure DB operations don't take too long
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	err = db.PingContext(ctx)
    if err != nil {
        db.Close()
        return nil, err
    }

    // return the connection pool (sql.DB)
    return db, nil
}