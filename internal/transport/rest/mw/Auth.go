package mw

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"shortener/errs"
	"shortener/internal/transport/rest"
)

type JwtOptions struct {
	Secret     []byte
	CookieName string
	UserIdKey  string
}

func InjectUserIdToCtx(options JwtOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := ExtractLog(r.Context(), "transport.rest.mw.InjectUserIdToCtx") //todo op and take in function

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

			userID, ok := (*claims)[options.UserIdKey].(float64)
			if !ok {
				log.Warn("Invalid user ID in token", slog.Any("token", token))
				return
			}

			ctx = context.WithValue(r.Context(), userIdKey{}, int64(userID))
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

func CheckAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := ExtractLog(r.Context(), "transport.rest.mw.InjectUserIdToCtx")
		_, ok := ExtractUserID(r.Context())

		if !ok {
			rest.Error(w, log, errs.E("", errors.New(""), errs.Unauthorized))
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
