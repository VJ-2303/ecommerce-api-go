package data

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int64     `json:"id"`
	Name        int64     `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Password    password  `json:"-"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}
