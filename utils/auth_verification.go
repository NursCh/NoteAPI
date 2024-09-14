package utils

import (
	"context"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		tokenString = tokenParts[1]
		claims, err := VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}
		//fmt.Println("Claims:", claims["userID"], claims["exp"])
		userIDfloat, ok := claims["userID"].(float64)
		userID := uint(userIDfloat)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
