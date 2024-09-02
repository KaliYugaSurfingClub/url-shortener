package handler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"link_shortener/internal/http/middlewares/mwLogger"
	"link_shortener/internal/lib/Api"
	"link_shortener/internal/storage"
	"log/slog"
	"net/http"
)

type urlGetter interface {
	GetURL(alias string) (string, error)
}

func Redirect(getter urlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mwLogger.GetCtxLog(r.Context(), "handlers.save.New")

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Warn("empty alias")
			render.JSON(w, r, Api.Error("empty alias"))
			return
		}

		log = log.With(slog.String("alias", alias))

		url, err := getter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info(storage.ErrURLNotFound.Error())
			render.JSON(w, r, Api.Error(storage.ErrURLNotFound.Error()))
			return
		}
		if err != nil {
			log.Error("getting url was failed", slog.String("error", err.Error()))
			render.JSON(w, r, Api.Error("getting url was failed"))
			return
		}

		log.Debug("do redirect", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
