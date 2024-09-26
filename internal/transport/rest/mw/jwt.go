package mw

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

// todo literal
const secretKey = "sasha"

func Jwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := ExtractLog(r.Context(), "transport.rest.Jwt")

		//todo literal
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			log.Error("Parse token:", ErrAttr(err))
			//todo response
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Error("Invalid token")
			//todo response
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		//todo literal
		username := (*claims)["username"].(string)
		userID := (*claims)["user_id"].(float64)

		ctx := context.WithValue(r.Context(), userNameKey{}, username)
		ctx = context.WithValue(ctx, userIdKey{}, int(userID))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type userNameKey struct{}

func ExtractUserName(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", errors.New("context is nil")
	}

	username, ok := ctx.Value(LogKey{}).(string)
	if !ok {
		return "", errors.New("could not extract username from context")
	}

	return username, nil
}

type userIdKey struct{}

func ExtractUserID(ctx context.Context) (int, error) {
	if ctx == nil {
		return 0, errors.New("context is nil")
	}

	id, ok := ctx.Value(userIdKey{}).(int)
	if !ok {
		return 0, errors.New("could not extract user id from context")
	}

	return id, nil
}
