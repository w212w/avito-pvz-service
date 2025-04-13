package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const RoleKey contextKey = "role"

func AuthMiddleware(secretKey string, requiredRole ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			tokenStr, err := extractTokenFromHeader(r)
			if err != nil {
				writeJSONError(w, http.StatusForbidden, "acces denied 1")
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})
			if err != nil || !token.Valid {
				writeJSONError(w, http.StatusForbidden, "acces denied 2")
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				writeJSONError(w, http.StatusForbidden, "acces denied 3")
				return
			}

			if len(requiredRole) > 0 {
				if !slices.Contains(requiredRole, role) {
					writeJSONError(w, http.StatusForbidden, "access denied from middleware")
					return
				}
			}

			ctx := context.WithValue(r.Context(), RoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}
	return parts[1], nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Message: message})
}

type errorResponse struct {
	Message string `json:"message"`
}
