package completeAdHandler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"shortener/internal/transport/rest/mw"
	"strconv"
)

type AdCompleter interface {
	CompleteAd(ctx context.Context, clickId int64)
}

type Handler struct {
	completer AdCompleter
}

func New(completer AdCompleter) *Handler {
	return &Handler{completer: completer}
}

func (h *Handler) Handler(_ http.ResponseWriter, r *http.Request) {
	log := mw.ExtractLog(r.Context(), "transport.rest.openShortenedHandler")

	clickId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Error(err.Error())
	}

	go h.completer.CompleteAd(r.Context(), clickId)
}
