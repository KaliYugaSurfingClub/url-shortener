package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"shortener/internal/transport/rest/mw"
)

type Rest struct {
}

func New(log *slog.Logger) *http.Server {
	r := chi.NewRouter()

	//todo good
	//r.Use(middleware.RealIP)

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mw.Logger(log))

	rest := &Rest{}

	r.Route("/link", func(r chi.Router) {
		r.Post("/", rest.LinkManager())
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return server
}

type LinkManager struct {
}

func (r *Rest) LinkManager() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.Rest.LinkManager")

		username, err := mw.ExtractUserName(r.Context())
		if err != nil {
			log.Error("could not extract user name", mw.ErrAttr(err))
			return
		}

		w.Write([]byte(username))
	}
}
