package shortLinkHandler

import (
	"context"
	"github.com/KaliYugaSurfingClub/errs/response"
	"github.com/go-chi/render"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
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

		defer r.Body.Close()

		req := &request{}
		if err := render.Decode(r, req); err != nil {
			//todo decode error
			return
		}

		shorted, err := shortener.Short(r.Context(), *req.ToModel(userId))
		if err != nil {
			response.Error(w, log, err)
			return
		}

		response.Ok(w, shorted)
	}
}
