package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/v1/users", a.createUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", a.showUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", a.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", a.deleteUserHandler)
	return a.recoverPanic(router)
}

