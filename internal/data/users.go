package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicatePhoneNumber = errors.New("duplicate phone number")

type User struct {
	ID          int64    `json:"id"`
	Name        int64    `json:"name"`
	PhoneNumber string   `json:"phone_number"`
	Password    password `json:"-"`
	Role        string   `json:"role"`
	CreatedAt   Time     `json:"created_at"`
}

type password struct {
	PlainText string
	hash      []byte
}

func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return err
	}
	p.PlainText = plainText
	p.hash = hash
	return nil
}

func (p *password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `INSERT into users (name,phone_number,password_hash)
					 RETURNING id,role,created_at
	`
	args := []any{user.Name, user.PhoneNumber, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return ErrDuplicatePhoneNumber
		}
		return err
	}
	return nil
}
