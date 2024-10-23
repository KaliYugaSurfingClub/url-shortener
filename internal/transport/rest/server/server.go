package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"shortener/internal/config"
	"shortener/internal/transport/rest/handler"
	"shortener/internal/transport/rest/mw"
	"time"
)

type Handlers struct {
	ShortLink     http.HandlerFunc
	GetUserLinks  http.HandlerFunc
	GetLinkClicks http.HandlerFunc
	OpenLink      http.HandlerFunc
	CompleteAd    http.HandlerFunc
}

func New(handlers Handlers, JWTOptions config.Auth, httpOptions config.HTTPServer, log *slog.Logger) *http.Server {
	router := chi.NewRouter()

	jwtOpt := mw.JwtOptions{
		Secret:     []byte(JWTOptions.JWTSecret),
		CookieName: JWTOptions.UserIdCookieKey,
		UserIdKey:  JWTOptions.UserIdJWTKey,
	}

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(mw.NewLogger(log))
	router.Use(mw.InjectUserIdToCtx(jwtOpt))

	//todo temporary
	router.Route("/users", func(r chi.Router) {
		r.Get("/login", handler.Login(jwtOpt, 24*time.Hour))
	})
	router.Get("/static/video", handler.StreamVideoHandler)
	//

	router.Route("/links", func(r chi.Router) {
		r.Use(mw.CheckAuth)
		r.Post("/", handlers.ShortLink)
		r.Get("/", handlers.GetUserLinks)
		r.Get("/{linkId}/clicks", handlers.GetLinkClicks)
	})

	router.Get("/open/{alias}", handlers.OpenLink)
	router.Post("/complete", handlers.CompleteAd)

	server := &http.Server{
		Handler:           router,
		Addr:              httpOptions.Address,
		WriteTimeout:      httpOptions.WriteTimeout,
		ReadHeaderTimeout: httpOptions.ReadHeaderTimeout,
		ReadTimeout:       httpOptions.ReadTimeout,
		IdleTimeout:       httpOptions.IdleTimeout,
	}

	return server
}
