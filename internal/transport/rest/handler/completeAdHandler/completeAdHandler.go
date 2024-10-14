package completeAdHandler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/response"
	"strconv"
)

type adCompleter interface {
	CompleteAd(ctx context.Context, clickId int64)
}

func New(completer adCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.openShortenedHandler")

		clickId, err := strconv.ParseInt(chi.URLParam(r, "clickId"), 10, 64)
		if err != nil {
			log.Error(err.Error())
		}

		go completer.CompleteAd(r.Context(), clickId)

		render.JSON(w, r, response.WithOk())
	}
}
