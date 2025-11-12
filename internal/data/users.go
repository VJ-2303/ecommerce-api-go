package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/VJ-2303/CityStars/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatePhoneNumber = errors.New("duplicate phone number")
	ErrUserNotFound         = errors.New("user not found")
)

type User struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	PhoneNumber string   `json:"phone_number"`
	Password    password `json:"-"`
	Role        string   `json:"role"`
	CreatedAt   Time     `json:"created_at"`
}

func ValidateUser(v *validator.Validator, u *User) {
	v.Check(len(u.Name) > 5, "name", "name must be provided and greater than 5 character")
	v.Check(len(u.Password.PlainText) > 8, "password", "password length must be greater than 8 characters")
	v.Check(len(u.Password.PlainText) < 72, "password", "password must be less than 72 characters")
	v.Check(validator.Matches(u.PhoneNumber, validator.PhoneNumberRegex), "phone_number", "provide an valid phone number")
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
					 VALUES($1,$2,$3)
					 RETURNING id,role,created_at
	`
	args := []any{user.Name, user.PhoneNumber, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
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

func (m UserModel) GetByPhoneNumber(phoneNumber string) (*User, error) {
	query := `
				SELECT id, name, phone_number, password_hash,role, created_at
				FROM users
				WHERE phone_number = $1
				    `
	var u User

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, phoneNumber).Scan(
		&u.ID,
		&u.Name,
		&u.PhoneNumber,
		&u.Password.hash,
		&u.Role,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
