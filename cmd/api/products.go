package main

import (
	"errors"
	"net/http"

	"github.com/VJ-2303/ecommerce-api-go/internal/data"
)

func (app *application) showProductHandler(w http.ResponseWriter, r *http.Request) {
	// Reas the id parameter using the helper function
	// and sending not found response if it returns an error
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	product, err := app.models.Products.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrProductNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"product": product})
}
