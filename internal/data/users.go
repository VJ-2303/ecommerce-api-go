package data

import (
	"database/sql"
)

type User struct {
	ID          int64    `json:"id"`
	Name        int64    `json:"name"`
	PhoneNumber string   `json:"phone_number"`
	Password    password `json:"-"`
	Role        string   `json:"role"`
	CreatedAt   Time     `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}
