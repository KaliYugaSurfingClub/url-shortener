package shortLinkHandler

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"net/http"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/response"
	"shortener/internal/utils/valkit"
	"time"
)

// todo add save with lifetime
type request struct {
	Original       string     `json:"original"`
	Alias          string     `json:"alias"`
	CustomName     string     `json:"customName"`
	ClicksToExpire *int64     `json:"clicksToExpire,omitempty"`
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
}

type LinkShortener interface {
	Short(ctx context.Context, link model.Link) (*model.Link, error)
}

func New(shortener LinkShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.Rest.ShortLink")

		link, err := linkFromRequest(r)
		if err != nil {
			log.Error("invalid request", mw.ErrAttr(err))
			render.JSON(w, r, response.NewError(err))
			return
		}

		shorted, err := shortener.Short(r.Context(), *link)
		if errors.Is(err, core.ErrAliasExists) { //todo
			render.JSON(w, r, response.NewError(core.ErrAliasExists))
			return
		}
		if errors.Is(err, core.ErrCustomNameExists) {
			render.JSON(w, r, response.NewError(core.ErrAliasExists))
			return
		}
		if err != nil {
			log.Error("cannot save link", mw.ErrAttr(err))
			render.JSON(w, r, response.NewInternalError())
			return
		}

		render.JSON(w, r, response.NewOk(response.LinkFromModel(shorted)))
	}
}

func linkFromRequest(r *http.Request) (*model.Link, error) {
	defer r.Body.Close()

	userId, _ := mw.ExtractUserID(r.Context())

	req := &request{}
	if err := render.DecodeJSON(r.Body, req); err != nil {
		return nil, err
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	return &model.Link{
		CreatedBy:      userId,
		Original:       req.Original,
		Alias:          req.Alias,
		CustomName:     req.CustomName,
		ExpirationDate: req.ExpirationDate,
		ClicksToExpire: req.ClicksToExpire,
	}, nil
}

func (r *request) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Original, validation.Required, is.URL),
		validation.Field(&r.Alias, validation.Length(1, 255)),      //todo learn from db
		validation.Field(&r.CustomName, validation.Length(1, 255)), //todo learn from db
		validation.Field(&r.ClicksToExpire, validation.By(valkit.IsPositive())),
		validation.Field(&r.ExpirationDate, validation.By(valkit.IsFutureDate())),
	)
}
