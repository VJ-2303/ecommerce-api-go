package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/VJ-2303/CityStars/internal/data"
	"github.com/VJ-2303/CityStars/internal/validator"
)

// CreateReportHandler allows authenticated users to create a new report
func (app *application) CreateReportHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDKey).(int64)

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Location    string `json:"location"`
		BeforeImage string `json:"before_image"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	report := &data.Report{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
		Location:    input.Location,
		BeforeImage: input.BeforeImage,
		Status:      "pending",
	}

	v := validator.New()
	if data.ValidateReport(v, report); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Reports.Insert(report)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"report": report})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// ListAllReportsHandler returns all reports with optional filtering (public endpoint)
func (app *application) ListAllReportsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	qs := r.URL.Query()

	limit := 50 // default limit
	if limitStr := qs.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := qs.Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	status := qs.Get("status")
	category := qs.Get("category")

	reports, err := app.models.Reports.GetAll(limit, offset, status, category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"reports": reports})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetReportHandler retrieves a single report by ID (public endpoint)
func (app *application) GetReportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	report, err := app.models.Reports.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrReportNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"report": report})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// UpdateReportStatusHandler allows admins to update report status and add after image
func (app *application) UpdateReportStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Status     string `json:"status"`
		AfterImage string `json:"after_image"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateReportUpdate(v, input.Status, input.AfterImage)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Reports.Update(id, input.Status, input.AfterImage)
	if err != nil {
		if errors.Is(err, data.ErrReportNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Retrieve the updated report
	report, err := app.models.Reports.Get(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"report": report})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetUserReportsHandler retrieves all reports created by the authenticated user
func (app *application) GetUserReportsHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDKey).(int64)

	// Parse query parameters for pagination
	qs := r.URL.Query()

	limit := 50 // default limit
	if limitStr := qs.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := qs.Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	reports, err := app.models.Reports.GetByUserID(userID, limit, offset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"reports": reports})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetReportStatsHandler returns statistics about reports (public endpoint)
func (app *application) GetReportStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := app.models.Reports.GetStats()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"stats": stats})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetLeaderboardHandler returns the top 10 users with most reports (public endpoint)
func (app *application) GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := app.models.Reports.GetLeaderboard()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"leaderboard": leaderboard})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
