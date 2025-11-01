package main

import (
	"net/http"
)

func (app *application) healtchCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	app.writeJSON(w, http.StatusOK, envelope{"server": data})
}
