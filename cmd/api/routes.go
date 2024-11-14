package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/books", a.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", a.displayBookHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", a.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", a.deleteBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books", a.listBooksHandler)

	router.HandlerFunc(http.MethodPost, "/v1/books/:id/reviews", a.createReviewHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id/reviews/:review_id", a.displayReviewHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id/reviews/:review_id", a.updateReviewHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id/reviews/:review_id", a.deleteReviewHandler)
	router.HandlerFunc(http.MethodGet, "/v1/reviews", a.listReviewsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id/reviews", a.listBookReviewsHandler)

	//return a.recoverPanic(router)
	return a.recoverPanic(a.rateLimit(router))
}


