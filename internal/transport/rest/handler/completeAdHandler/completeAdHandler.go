package completeAdHandler

import (
	"context"
	"net/http"
	"shortener/internal/transport/rest"
	"shortener/internal/transport/rest/mw"
)

//todo send original

type adCompleter interface {
	CompleteAd(ctx context.Context, clickId int64)
}

type request struct {
	ClickId int64 `json:"clickId"`
}

func New(completer adCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.openShortenedHandler")

		req := new(request)
		if err := rest.DecodeJSON(req, r); err != nil {
			rest.Error(w, log, err)
			return
		}

		go completer.CompleteAd(r.Context(), req.ClickId)
	}
}
