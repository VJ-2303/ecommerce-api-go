package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/VJ-2303/CityStars/internal/validator"
)

var (
	ErrReportNotFound = errors.New("report not found")
)

// Report represents a problem report submitted by a citizen
type Report struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Location    string `json:"location"`
	BeforeImage string `json:"before_image"`
	AfterImage  string `json:"after_image,omitempty"`
	Status      string `json:"status"`
	CreatedAt   Time   `json:"created_at"`
	UpdatedAt   Time   `json:"updated_at"`
	CompletedAt *Time  `json:"completed_at,omitempty"`
	UserName    string `json:"user_name,omitempty"`
}

// ValidateReport validates the report data
func ValidateReport(v *validator.Validator, r *Report) {
	v.Check(len(r.Title) > 0, "title", "title must be provided")
	v.Check(len(r.Title) <= 200, "title", "title must not be more than 200 characters")
	v.Check(len(r.Description) > 0, "description", "description must be provided")
	v.Check(len(r.Description) <= 2000, "description", "description must not be more than 2000 characters")
	v.Check(len(r.Category) > 0, "category", "category must be provided")
	v.Check(validator.PermittedValue(r.Category, "pothole", "streetlight", "water", "garbage", "road", "other"), "category", "invalid category")
	v.Check(len(r.Location) > 0, "location", "location must be provided")
	v.Check(len(r.Location) <= 500, "location", "location must not be more than 500 characters")
	v.Check(len(r.BeforeImage) > 0, "before_image", "before image must be provided")
}

// ValidateReportUpdate validates report update data
func ValidateReportUpdate(v *validator.Validator, status, afterImage string) {
	v.Check(len(status) > 0, "status", "status must be provided")
	v.Check(validator.PermittedValue(status, "pending", "in-progress", "completed", "rejected"), "status", "invalid status")
	if status == "completed" {
		v.Check(len(afterImage) > 0, "after_image", "after image is required when marking as completed")
	}
}

// ReportModel wraps the database connection
type ReportModel struct {
	DB *sql.DB
}

// Insert creates a new report in the database
func (m ReportModel) Insert(report *Report) error {
	query := `
		INSERT INTO reports (user_id, title, description, category, location, before_image, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	args := []any{
		report.UserID,
		report.Title,
		report.Description,
		report.Category,
		report.Location,
		report.BeforeImage,
		"pending",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&report.ID,
		&report.CreatedAt,
		&report.UpdatedAt,
	)
}

// Get retrieves a single report by ID with user information
func (m ReportModel) Get(id int64) (*Report, error) {
	query := `
		SELECT r.id, r.user_id, r.title, r.description, r.category, r.location, 
		       r.before_image, r.after_image, r.status, r.created_at, r.updated_at, 
		       r.completed_at, u.name as user_name
		FROM reports r
		INNER JOIN users u ON r.user_id = u.id
		WHERE r.id = $1
	`

	var report Report
	var completedAt sql.NullTime

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&report.ID,
		&report.UserID,
		&report.Title,
		&report.Description,
		&report.Category,
		&report.Location,
		&report.BeforeImage,
		&report.AfterImage,
		&report.Status,
		&report.CreatedAt,
		&report.UpdatedAt,
		&completedAt,
		&report.UserName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReportNotFound
		}
		return nil, err
	}

	if completedAt.Valid {
		t := Time(completedAt.Time)
		report.CompletedAt = &t
	}

	return &report, nil
}

// GetAll retrieves all reports with pagination
func (m ReportModel) GetAll(limit, offset int, status, category string) ([]*Report, error) {
	query := `
		SELECT r.id, r.user_id, r.title, r.description, r.category, r.location, 
		       r.before_image, r.after_image, r.status, r.created_at, r.updated_at, 
		       r.completed_at, u.name as user_name
		FROM reports r
		INNER JOIN users u ON r.user_id = u.id
		WHERE ($3 = '' OR r.status = $3)
		  AND ($4 = '' OR r.category = $4)
		ORDER BY r.created_at DESC
		LIMIT $1 OFFSET $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, limit, offset, status, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []*Report{}

	for rows.Next() {
		var report Report
		var completedAt sql.NullTime

		err := rows.Scan(
			&report.ID,
			&report.UserID,
			&report.Title,
			&report.Description,
			&report.Category,
			&report.Location,
			&report.BeforeImage,
			&report.AfterImage,
			&report.Status,
			&report.CreatedAt,
			&report.UpdatedAt,
			&completedAt,
			&report.UserName,
		)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			t := Time(completedAt.Time)
			report.CompletedAt = &t
		}

		reports = append(reports, &report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

// GetByUserID retrieves all reports by a specific user
func (m ReportModel) GetByUserID(userID int64, limit, offset int) ([]*Report, error) {
	query := `
		SELECT r.id, r.user_id, r.title, r.description, r.category, r.location, 
		       r.before_image, r.after_image, r.status, r.created_at, r.updated_at, 
		       r.completed_at, u.name as user_name
		FROM reports r
		INNER JOIN users u ON r.user_id = u.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []*Report{}

	for rows.Next() {
		var report Report
		var completedAt sql.NullTime

		err := rows.Scan(
			&report.ID,
			&report.UserID,
			&report.Title,
			&report.Description,
			&report.Category,
			&report.Location,
			&report.BeforeImage,
			&report.AfterImage,
			&report.Status,
			&report.CreatedAt,
			&report.UpdatedAt,
			&completedAt,
			&report.UserName,
		)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			t := Time(completedAt.Time)
			report.CompletedAt = &t
		}

		reports = append(reports, &report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

// Update updates a report's status and after image (only admin can do this)
func (m ReportModel) Update(id int64, status, afterImage string) error {
	query := `
		UPDATE reports
		SET status = $1, 
		    after_image = $2, 
		    updated_at = NOW(),
		    completed_at = CASE WHEN $1 = 'completed' THEN NOW() ELSE completed_at END
		WHERE id = $3
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	var returnedID int64
	err := m.DB.QueryRowContext(ctx, query, status, afterImage, id).Scan(&returnedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrReportNotFound
		}
		return err
	}

	return nil
}

// ReportStats represents the statistics of reports
type ReportStats struct {
	TotalReports      int `json:"total_reports"`
	PendingReports    int `json:"pending_reports"`
	InProgressReports int `json:"in_progress_reports"`
	CompletedReports  int `json:"completed_reports"`
	RejectedReports   int `json:"rejected_reports"`
}

// GetStats retrieves report statistics
func (m ReportModel) GetStats() (*ReportStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'in-progress' THEN 1 END) as in_progress,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected
		FROM reports
	`

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	var stats ReportStats
	err := m.DB.QueryRowContext(ctx, query).Scan(
		&stats.TotalReports,
		&stats.PendingReports,
		&stats.InProgressReports,
		&stats.CompletedReports,
		&stats.RejectedReports,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// LeaderboardEntry represents a user in the leaderboard
type LeaderboardEntry struct {
	UserID      int64  `json:"user_id"`
	UserName    string `json:"user_name"`
	ReportCount int    `json:"report_count"`
	PhoneNumber string `json:"phone_number"`
}

// GetLeaderboard retrieves top 10 users with most reports
func (m ReportModel) GetLeaderboard() ([]*LeaderboardEntry, error) {
	query := `
		SELECT u.id, u.name, u.phone_number, COUNT(r.id) as report_count
		FROM users u
		INNER JOIN reports r ON u.id = r.user_id
		GROUP BY u.id, u.name, u.phone_number
		ORDER BY report_count DESC
		LIMIT 10
	`

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaderboard := []*LeaderboardEntry{}

	for rows.Next() {
		var entry LeaderboardEntry
		err := rows.Scan(
			&entry.UserID,
			&entry.UserName,
			&entry.PhoneNumber,
			&entry.ReportCount,
		)
		if err != nil {
			return nil, err
		}

		leaderboard = append(leaderboard, &entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return leaderboard, nil
}
