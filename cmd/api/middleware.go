package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userIDKey   = contextKey("userID")
	userRoleKey = contextKey("role")
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request received", "method", r.Method, "uri", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		tokenString := headerParts[1]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return []byte(app.config.jwtSecret), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				app.expiredTokenResponse(w, r)
			} else {
				app.invalidAuthenticationTokenResponse(w, r)
			}
			return
		}
		if !token.Valid {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		sub, err := claims.GetSubject()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		userID, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		role := claims["role"].(string)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx2 := context.WithValue(ctx, userRoleKey, role)

		newReq := r.WithContext(ctx2)

		next.ServeHTTP(w, newReq)
	})
}

func (app *application) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(userRoleKey).(string)
		if role != "admin" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
