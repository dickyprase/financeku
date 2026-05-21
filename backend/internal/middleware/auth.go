package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/financeku/backend/internal/config"
	"github.com/financeku/backend/pkg/response"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func RequireAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "Missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(w, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				response.Unauthorized(w, "Invalid or expired token")
				return
			}

			tokenType, _ := claims["type"].(string)
			if tokenType != "access" {
				response.Unauthorized(w, "Invalid token type")
				return
			}

			userID, _ := claims["user_id"].(string)
			role, _ := claims["role"].(string)

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, RoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(RoleKey).(string)
		if !ok || role != "admin" {
			response.Forbidden(w, "Admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) string {
	userID, _ := r.Context().Value(UserIDKey).(string)
	return userID
}

func GetUserRole(r *http.Request) string {
	role, _ := r.Context().Value(RoleKey).(string)
	return role
}
