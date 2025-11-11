package main

import "net/http"

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request received", "method", r.Method, "uri", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
