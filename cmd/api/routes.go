package main

import (
	//"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies)routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler) 
	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	router.HandlerFunc(http.MethodPatch,"/v1/comments/:id", a.updateCommentHandler)
	router.HandlerFunc(http.MethodDelete,"/v1/comments/:id", a.deleteCommentHandler)
	router.HandlerFunc(http.MethodGet,"/v1/comments", a.listCommentsHandler)

	return a.recoverPanic(router)
}