package completeAdHandler

import (
	"context"
	"github.com/KaliYugaSurfingClub/pkg/mw"
	"link-service/internal/transport/rest"
	"net/http"
)

type adCompleter interface {
	CompleteAd(ctx context.Context, clickId int64) (string, error)
}

type request struct {
	ClickId int64 `json:"click_id"`
}

type response struct {
	Original string `json:"original"`
}

func New(completer adCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.openShortenedHandler")

		req := new(request)
		if err := rest.DecodeJSON(req, r); err != nil {
			rest.Error(w, log, err)
			return
		}

		original, err := completer.CompleteAd(r.Context(), req.ClickId)
		if err != nil {
			rest.Error(w, log, err)
			return
		}

		rest.Ok(w, response{
			Original: original,
		})
	}
}
