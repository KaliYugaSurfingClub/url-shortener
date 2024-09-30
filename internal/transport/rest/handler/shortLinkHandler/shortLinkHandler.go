package shortLinkHandler

import (
	"context"
	"github.com/go-chi/render"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/response"
	"time"
)

type request struct {
	Original       string     `json:"original"`
	Alias          string     `json:"alias"`
	CustomName     string     `json:"customName"`
	ClicksToExpire *int64     `json:"clicksToExpire,omitempty"`
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
}

type data struct {
	Alias      string `json:"alias"`
	CustomName string `json:"customName"`
}

type LinkShortener interface {
	Short(ctx context.Context, link model.Link) (*model.Link, error)
}

func New(shortener LinkShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.Rest.ShortLink")

		id, _ := mw.ExtractUserID(r.Context())

		//todo validate and scan into model.Link in one func
		req := &request{}
		if err := render.DecodeJSON(r.Body, req); err != nil {
			log.Error("cannot decode request", mw.ErrAttr(err))
			render.JSON(w, r, response.NewError(err))
			return
		}

		link := model.Link{
			CreatedBy:      id,
			Original:       req.Original,
			Alias:          req.Alias,
			CustomName:     req.CustomName,
			ExpirationDate: req.ExpirationDate,
			ClicksToExpire: req.ClicksToExpire,
		}

		shorted, err := shortener.Short(r.Context(), link)
		if err != nil {
			log.Error("cannot save link", mw.ErrAttr(err))
			render.JSON(w, r, response.NewError(err))
			return
		}

		render.JSON(w, r, response.NewOk(data{
			Alias:      shorted.Alias,
			CustomName: shorted.CustomName,
		}))
	}
}
