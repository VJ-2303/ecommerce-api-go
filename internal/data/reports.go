package data

import (
	"context"
	"database/sql"
	"time"
)

type Report struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Status      string    `json:"status"`
	BeforeImage string    `json:"before_img"`
	AfterImage  string    `json:"after_img"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type ReportModel struct {
	DB *sql.DB
}

func (m ReportModel) Create(r *Report) error {
	query := `INSERT INTO reports
		    (title,description,address,status,before_img,after_img,created_by)
			VALUES($1,$2,$3,$4,$5,$6,$7)
			RETURNING id,created_at`

	args := []any{r.Title, r.Description, r.Address, r.Status, r.BeforeImage, r.AfterImage, r.CreatedBy}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&r.ID, r.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
