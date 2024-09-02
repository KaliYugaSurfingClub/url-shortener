package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"link_shortener/internal/lib/response"
	"link_shortener/internal/storage"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, getter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		alias := chi.URLParam(r, "alias")

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if alias == "" {
			log.Warn("empty alias")
			render.JSON(w, r, response.Error("empty alias"))
			return
		}

		log = log.With(slog.String("alias", alias))

		url, err := getter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info(storage.ErrURLNotFound.Error())
			render.JSON(w, r, response.Error(storage.ErrURLNotFound.Error()))
			return
		}
		if err != nil {
			log.Error("getting url was failed", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("getting url was failed"))
			return
		}

		log.Debug("do redirect", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
