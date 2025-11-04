package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new Router from
	router := httprouter.New()

	// Set the NotFound field of router to our custom notfoundResponse handler
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Healthcheck route
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healtchCheckHandler)

	// Products specific routes
	router.HandlerFunc(http.MethodGet, "/v1/product/:id", app.showProductHandler)
	router.HandlerFunc(http.MethodPost, "/v1/product", app.createProductHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/product/:id", app.updateProductHandler)

	// Return the router
	return router
}
