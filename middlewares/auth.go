package middlewares

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"storex/utils"
	"strings"
)

type contextKey string

const (
	UserIDKey contextKey = "userId"
	RoleKey   contextKey = "userRole"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logrus.Info("no Authorization header")
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			logrus.Info("invalid Authorization header")
			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseAccessToken(tokenString)
		if err != nil {
			logrus.Info("invalid token")
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			logrus.Info("invalid user id")
			http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
			return
		}

		// Inject into request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
