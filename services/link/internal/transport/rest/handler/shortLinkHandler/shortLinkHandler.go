package shortLinkHandler

import (
	"context"
	"github.com/KaliYugaSurfingClub/pkg/mw"
	"link-service/internal/core/model"
	"link-service/internal/transport/rest"
	"net/http"
)

type LinkShortener interface {
	Short(ctx context.Context, link model.Link) (*model.Link, error)
}

type request struct {
	Original   string `json:"original"`
	Alias      string `json:"alias"`
	CustomName string `json:"customName"`
}

func (r *request) ToModel(userId int64) *model.Link {
	return &model.Link{
		CreatedBy:  userId,
		Original:   r.Original,
		Alias:      r.Alias,
		CustomName: r.CustomName,
	}
}

func New(shortener LinkShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.Rest.ShortLink")
		userId, _ := mw.ExtractUserID(r.Context())

		req := &request{}
		if err := rest.DecodeJSON(req, r); err != nil {
			rest.Error(w, log, err)
			return
		}

		shorted, err := shortener.Short(r.Context(), *req.ToModel(userId))
		if err != nil {
			rest.Error(w, log, err)
			return
		}

		rest.Ok(w, shorted)
	}
}
