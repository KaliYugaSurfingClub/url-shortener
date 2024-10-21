package completeAdHandler

import (
	"context"
	"github.com/go-chi/render"
	"net/http"
	"shortener/internal/transport/rest/mw"
)

type adCompleter interface {
	CompleteAd(ctx context.Context, clickId int64)
}

//todo send original

type request struct {
	ClickId int64 `json:"clickId"`
}

func New(completer adCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = mw.ExtractLog(r.Context(), "transport.rest.openShortenedHandler")

		req := new(request)
		if err := render.Decode(r, req); err != nil {
			//todo decode error
			return
		}

		go completer.CompleteAd(r.Context(), req.ClickId)
	}
}
