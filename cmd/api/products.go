package main

import (
	"errors"
	"net/http"

	"github.com/VJ-2303/ecommerce-api-go/internal/data"
	"github.com/VJ-2303/ecommerce-api-go/internal/validator"
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

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		Price          int64  `json:"price"`
		StockAvailable int    `json:"stock_available"`
		ImageURL       string `json:"image_url"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	product := &data.Product{
		Name:           input.Name,
		Description:    input.Description,
		Price:          input.Price,
		StockAvailable: input.StockAvailable,
		ImageURL:       input.ImageURL,
	}
	v := validator.New()

	if data.ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Products.Insert(product)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"product": product})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
