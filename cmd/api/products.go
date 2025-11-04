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

	// Query the product using id
	product, err := app.models.Products.Get(id)
	if err != nil {
		// Check for our custom error type and send not found response
		// otherwise sent and generic serverErrorResponse
		if errors.Is(err, data.ErrProductNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the whole product struct in the response body
	err = app.writeJSON(w, http.StatusOK, envelope{"product": product})
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	// An temprory struct for decoding the JSON Body from the request
	var input struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		Price          int64  `json:"price"`
		StockAvailable int    `json:"stock_available"`
		ImageURL       string `json:"image_url"`
	}

	// decoding the JSON body into the anonymous struct by providing address
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Initiate an product and map the fields from input struct
	product := &data.Product{
		Name:           input.Name,
		Description:    input.Description,
		Price:          input.Price,
		StockAvailable: input.StockAvailable,
		ImageURL:       input.ImageURL,
	}
	// New Validator initiaded
	v := validator.New()

	// Validatation performed on the newly created product struct
	if data.ValidateProduct(v, product); !v.Valid() {
		// Sent custom error response if validation failed along with
		// Validatator.Errors map in the response
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Calling the Insert method on Products Model to
	// create an new product, it Creates new row in the products table
	// and Fill the id, created_at, updated_at to the product struct
	err = app.models.Products.Insert(product)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the complete product struct to the client
	// Along the information they provided and along with system generated
	// id, created_at and updated_at
	err = app.writeJSON(w, http.StatusCreated, envelope{"product": product})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	p, err := app.models.Products.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrProductNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Name           *string `json:"name"`
		Description    *string `json:"description"`
		Price          *int64  `json:"price"`
		StockAvailable *int    `json:"stock_available"`
		ImageURL       *string `json:"image_url"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Description != nil {
		p.Description = *input.Description
	}
	if input.Price != nil {
		p.Price = *input.Price
	}
	if input.StockAvailable != nil {
		p.StockAvailable = *input.StockAvailable
	}
	if input.ImageURL != nil {
		p.ImageURL = *input.ImageURL
	}
	v := validator.New()
	if data.ValidateProduct(v, p); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Products.Update(p)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"product": p})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
