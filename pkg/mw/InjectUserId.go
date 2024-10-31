package mw

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
)

type userIdKey struct{}

type InjectUserOptions struct {
	Secret     []byte
	CookieName string
	JWTKey     string
}

// InjectUserIdToCtx use mw.Logger
func InjectUserIdToCtx(options InjectUserOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := ExtractLog(r.Context(), "transport.rest.mw.InjectUserIdToCtx")
			if log == nil {
				panic("log is nil")
			}

			ctx := r.Context()

			defer func() {
				next.ServeHTTP(w, r.WithContext(ctx))
			}()

			cookie, err := r.Cookie(options.CookieName) //if userId are necessary, err will happen be in CheckAuth
			if err != nil {
				return
			}

			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return options.Secret, nil
			})

			if err != nil {
				log.Warn("Parse token:", ErrAttr(err))
				return
			}

			if !token.Valid {
				log.Warn("Invalid token", slog.Any("token", token))
				return
			}

			userID, ok := (*claims)[options.JWTKey].(float64)
			if !ok {
				log.Warn("Invalid user ID in token", slog.Any("token", token))
				return
			}

			ctx = context.WithValue(r.Context(), userIdKey{}, int64(userID))
		}

		return http.HandlerFunc(fn)
	}
}

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
