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
	err := app.writeJSON(w, http.StatusOK, envelope{"server": data})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
