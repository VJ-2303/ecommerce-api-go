package main

import (
	"errors"
	"net/http"

	"github.com/VJ-2303/ecommerce-api-go/internal/data"
	"github.com/VJ-2303/ecommerce-api-go/internal/validator"
)

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		PlainText   string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	app.logger.Debug(input.Name, input.PhoneNumber, input.PlainText)
	user := &data.User{
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Role:        "user",
	}
	err = user.Password.Set(input.PlainText)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicatePhoneNumber) {
			v.AddError("phone_number", "This phone number is already exists")
			app.failedValidationResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
