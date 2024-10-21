package mw

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

//todo refactor

type JwtOptions struct {
	Secret     []byte
	CookieName string
	UserIdKey  string
}

func CheckAuth(options JwtOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := ExtractLog(r.Context(), "transport.rest.mw.Jwt")

			cookie, err := r.Cookie(options.CookieName)
			if errors.Is(err, http.ErrNoCookie) {
				log.Error("Parse cookie:", ErrAttr(err))
				//todo
				return
			}

			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return options.Secret, nil
			})

			if err != nil {
				log.Error("Parse token:", ErrAttr(err))
				//todo
				return
			}

			if !token.Valid {
				log.Error("Invalid token")
				//todo
				return
			}

			userID := (*claims)[options.UserIdKey].(float64)

			ctx := context.WithValue(r.Context(), userIdKey{}, int64(userID))

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

type userIdKey struct{}

func ExtractUserID(ctx context.Context) (int64, bool) {
	if ctx == nil {
		return 0, false
	}

	id, ok := ctx.Value(userIdKey{}).(int64)
	if !ok {
		return 0, false
	}

	return id, true
}
