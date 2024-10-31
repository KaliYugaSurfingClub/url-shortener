package server

import (
	"github.com/KaliYugaSurfingClub/pkg/mw"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"link-service/internal/config"
	"link-service/internal/transport/rest/handler"
	"log/slog"
	"net/http"
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

	injectUserOpt := mw.InjectUserOptions{
		Secret:     []byte(JWTOptions.JWTSecret),
		CookieName: JWTOptions.CookieName,
		JWTKey:     JWTOptions.JWTKey,
	}

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(mw.Logger(log))
	router.Use(mw.InjectUserIdToCtx(injectUserOpt))

	//todo temporary
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
		Addr:              httpOptions.IP + ":" + httpOptions.Port,
		WriteTimeout:      httpOptions.WriteTimeout,
		ReadHeaderTimeout: httpOptions.ReadHeaderTimeout,
		ReadTimeout:       httpOptions.ReadTimeout,
		IdleTimeout:       httpOptions.IdleTimeout,
	}

	return server
}
