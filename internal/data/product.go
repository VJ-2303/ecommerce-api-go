package data

import (
	"database/sql"
	"time"
)

type Product struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          int64     `json:"price"`
	StockAvailable int       `json:"stock_available"`
	ImageURL       string    `json:"image_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProductModel struct {
	DB *sql.DB
}
