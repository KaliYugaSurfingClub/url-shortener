package shortLinkHandler

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/response"
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
			log.Info("cannot decode body", mw.ErrAttr(err))
			render.JSON(w, r, response.WithError(err)) //todo
			return
		}

		shorted, err := shortener.Short(r.Context(), *req.ToModel(userId))
		if errors.Is(err, core.ErrAliasExists) {
			render.JSON(w, r, response.WithError(core.ErrAliasExists))
			return
		}
		if errors.Is(err, core.ErrCustomNameExists) {
			render.JSON(w, r, response.WithError(core.ErrCustomNameExists))
			return
		}
		if err != nil {
			log.Error("cannot save link", mw.ErrAttr(err))
			render.JSON(w, r, response.WithInternalError())
			return
		}

		render.JSON(w, r, response.WithData(response.LinkFromModel(shorted)))
	}
}
