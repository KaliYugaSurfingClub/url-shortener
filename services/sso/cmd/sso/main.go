package main

import (
	"github.com/KaliYugaSurfingClub/pkg/mw"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"os"
	"sso/internal/config"
	"time"
)

func Login(config config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "Login")

		claims := jwt.MapClaims{
			config.JWT.UserIdKey: 1, //todo
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(config.JWT.Secret))
		if err != nil { //todo
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Unable to sign token", mw.ErrAttr(err))
			return
		}

		cookie := &http.Cookie{
			Name:     config.Cookie.Name,
			Value:    tokenString,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(config.Cookie.Lifetime),
		}

		http.SetCookie(w, cookie)
	}
}

//todo

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(mw.Logger(log))
	router.Use(mw.InjectUserIdToCtx(mw.InjectUserOptions{
		Secret:     []byte(cfg.JWT.Secret),
		CookieName: cfg.Cookie.Name,
		JWTKey:     cfg.JWT.UserIdKey,
	}))

	router.Route("/users", func(r chi.Router) {
		r.Get("/login", Login(*cfg))
	})

	address := cfg.HTTP.IP + ":" + cfg.HTTP.Port
	log.Info("Starting server on " + address)

	log.Error(http.ListenAndServe(address, router).Error())
}
