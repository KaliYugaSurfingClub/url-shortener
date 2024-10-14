package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"shortener/internal/transport/rest/mw"
	"time"
)

//todo update

func Login(jwtOpt mw.JwtOptions, lifetime time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.Login")

		claims := jwt.MapClaims{
			jwtOpt.UserIdKey: 1,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtOpt.Secret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Unable to sign token", mw.ErrAttr(err))
			return
		}

		cookie := &http.Cookie{
			Name:     jwtOpt.CookieName,
			Value:    tokenString,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(lifetime),
		}

		http.SetCookie(w, cookie)
	}
}
