package handler

import (
	"context"
	"github.com/go-chi/render"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest"
	"shortener/internal/transport/rest/mw"
	"time"
)

type Shortener interface {
	Short(ctx context.Context, link model.Link) (*model.Link, error)
}

func Short(shortener Shortener) http.HandlerFunc {
	type request struct {
		Original       string     `json:"original"`
		Alias          string     `json:"alias"`
		CustomName     string     `json:"customName"`
		ClicksToExpire *int64     `json:"clicksToExpire,omitempty"`
		ExpirationDate *time.Time `json:"expirationDate,omitempty"`
	}

	type response struct {
		Alias      string `json:"alias"`
		CustomName string `json:"customName"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.Rest.LinkManager")

		id, _ := mw.ExtractUserID(r.Context())

		req := &request{}
		if err := render.DecodeJSON(r.Body, req); err != nil {
			log.Error("cannot decode request", mw.ErrAttr(err))
			render.JSON(w, r, rest.NewErrorResponse(err))
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
			render.JSON(w, r, rest.NewErrorResponse(err))
			return
		}

		render.JSON(w, r, &response{
			Alias:      shorted.Alias,
			CustomName: shorted.CustomName,
		})
	}
}
