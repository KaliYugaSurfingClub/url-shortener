package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"link_shortener/internal/lib/Api"
	"link_shortener/internal/lib/sl"
	"link_shortener/internal/storage"
	"link_shortener/internal/transport/middlewares/mwLogger"
	"log/slog"
	"net/http"
)

type urlGetter interface {
	GetURL(alias string) (string, error)
}

func New(getter urlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mwLogger.GetCtxLog(r.Context(), "handlers.save.New")

		alias := chi.URLParam(r, "alias")

		log = log.With(slog.String("alias", alias))

		url, err := getter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info(err.Error())
			render.JSON(w, r, Api.Error(err.Error()))
			return
		}
		if err != nil {
			log.Error("getting url was failed", sl.ErrorAttr(err))
			render.JSON(w, r, Api.Error("getting url was failed"))
			return
		}

		log.Debug("do redirect", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
