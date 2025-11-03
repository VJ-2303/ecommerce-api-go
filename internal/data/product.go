package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Custom error for telling is not found in the DB
var ErrProductNotFound = errors.New("Product not found")

// Product represents and single row in a Product table in the DB
type Product struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          int64     `json:"price"`
	StockAvailable int       `json:"stock_available"`
	ImageURL       *string   `json:"image_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ProductModel encloses and sql connections
// Methods such as Insert, Update, Get, Delete will be implemented against it
type ProductModel struct {
	DB *sql.DB
}

func (m ProductModel) Insert(product *Product) error {
	query := `INSERT INTO products
					(name,description,price,stock_available,image_url)
					VALUES ($1,$2,$3,$4,$5)
					RETURNING id, created_at, updated_at`

	args := []any{product.Name, product.Description, product.Price, product.StockAvailable, product.ImageURL}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

// Get method will an Product struct populated by the data
// from the DB, it takes the primary key 'id' as a argument
// and query the DB using the ID, it returns an Custom Error
// when the Product is not found
func (m ProductModel) Get(id int64) (*Product, error) {
	// If the id < 1 we return immediatly
	// because DB primary keys always starts from 1
	if id < 1 {
		return nil, ErrProductNotFound
	}
	query := `SELECT id,name,description,price,stock_available,image_url,created_at,updated_at
			 FROM products
			 WHERE id = $1`

	var p Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.StockAvailable,
		&p.ImageURL,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		// Return our custom Error when we sql.ErrNoRows from the DB
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		} else {
			return nil, err
		}
	}
	return &p, nil
}
