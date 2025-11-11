package main

import (
	"errors"
	"net/http"
	"time"

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

func (app *application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	v.Check(validator.Matches(input.PhoneNumber, validator.PhoneNumberRegex), "phone_number", "provide an valid phone number")
	v.Check(len(input.PhoneNumber) > 8, "password", "password must be atleast 8 characters long")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.Users.GetByPhoneNumber(input.PhoneNumber)
	if err != nil {
		if errors.Is(err, data.ErrUserNotFound) {
			app.authenticationErrorResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.authenticationErrorResponse(w, r)
		return
	}
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, "authentication", app.config.jwtSecret, user.Role)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"auth_token": token})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
