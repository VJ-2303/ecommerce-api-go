package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	PlainText string    `json:"token"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int64, expiry time.Duration, Scope, secretKey, role string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(expiry),
		Scope:  Scope,
	}
	claims := jwt.MapClaims{}
	claims["exp"] = token.Expiry.Unix()
	claims["iat"] = time.Now()
	claims["scope"] = Scope
	claims["sub"] = fmt.Sprintf("%d", userID)
	claims["role"] = role

	jwtString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := jwtString.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	token.PlainText = signedString
	return token, nil
}
