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

	// user routes
	router.HandlerFunc(http.MethodPost, "/v1/user/register", app.CreateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/user/login", app.LoginUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/user/me", app.authenticate(app.userProfileHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/me", app.authenticate(app.requireAdmin(app.AdminProfileHandler)))

	// Return the router
	return app.logRequest(router)
}
