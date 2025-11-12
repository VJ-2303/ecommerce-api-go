package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// Initialize a new chi router
	router := chi.NewRouter()

	// Use chi's built-in middleware
	router.Use(middleware.Recoverer)

	// Set custom handlers for 404 and 405
	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	// Healthcheck route
	router.Get("/v1/healthcheck", app.healtchCheckHandler)

	// User routes
	router.Post("/v1/user/register", app.CreateUserHandler)
	router.Post("/v1/user/login", app.LoginUserHandler)
	router.Get("/v1/user/me", app.authenticate(app.userProfileHandler))
	router.Get("/v1/user/reports", app.authenticate(app.GetUserReportsHandler))
	router.Get("/v1/admin/me", app.authenticate(app.requireAdmin(app.AdminProfileHandler)))

	// Report routes (Public - anyone can view)
	router.Get("/v1/reports", app.ListAllReportsHandler)
	router.Get("/v1/reports/stats", app.GetReportStatsHandler)
	router.Get("/v1/reports/{id}", app.GetReportHandler)
	router.Get("/v1/leaderboard", app.GetLeaderboardHandler)

	// Report routes (Authenticated users - create)
	router.Post("/v1/reports", app.authenticate(app.CreateReportHandler))

	// Admin routes - update report status
	router.Patch("/v1/reports/{id}", app.authenticate(app.requireAdmin(app.UpdateReportStatusHandler)))

	// Return the router with logging
	return app.logRequest(router)
}
