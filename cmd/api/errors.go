package main

import (
	"net/http"
)

// HanlerFunction for sending notFound error message
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "this resourse is not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// Handler for sending methodNotAllowed error message
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := "This method is not allowed on this specific route"
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// logError used to log any error happened in a specific request
// along with its method and request uri
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// errorResponse sends an generic error message to the client enclosed by the envelope
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := app.writeJSON(w, status, envelope{"error": message})
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// this response is sented when the server had any unknown errors
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "The server encountered and error and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// this response it sended when the request body is badly formed
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse is sented when the Input have validation errors,
// it takes the Errors map as argument and write the whole map inside the response body
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
