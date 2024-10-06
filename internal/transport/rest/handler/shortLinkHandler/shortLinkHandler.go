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

type LinkShortener interface {
	Short(ctx context.Context, link model.Link) (*model.Link, error)
}

type Handler struct {
	shortener        LinkShortener
	OriginalMaxLen   int
	AliasMaxLen      int
	CustomNameMaxLen int
}

func New(shortener LinkShortener, OriginalMaxLen, AliasMaxLen, CustomLenMaxLen int) *Handler {
	return &Handler{
		shortener:        shortener,
		OriginalMaxLen:   OriginalMaxLen,
		AliasMaxLen:      AliasMaxLen,
		CustomNameMaxLen: CustomLenMaxLen,
	}
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log := mw.ExtractLog(r.Context(), "transport.Rest.ShortLink")
	userId, _ := mw.ExtractUserID(r.Context())

	req := &request{}
	if err := render.Decode(r, req); err != nil {
		log.Info("cannot decode body", mw.ErrAttr(err))
		render.JSON(w, r, response.WithError(err)) //todo
		return
	}

	if errs := h.validateRequest(req); errs != nil {
		log.Info("invalid request", mw.ErrAttr(errs))
		render.JSON(w, r, response.WithValidationErrors(errs))
		return
	}

	shorted, err := h.shortener.Short(r.Context(), *req.ToModel(userId))
	if errors.Is(err, core.ErrAliasExists) {
		render.JSON(w, r, response.WithError(core.ErrAliasExists))
		return
	}
	if errors.Is(err, core.ErrCustomNameExists) {
		render.JSON(w, r, response.WithError(core.ErrAliasExists))
		return
	}
	if err != nil {
		log.Error("cannot save link", mw.ErrAttr(err))
		render.JSON(w, r, response.WithInternalError())
		return
	}

	render.JSON(w, r, response.WithOk(response.LinkFromModel(shorted)))
}

// todo add save with lifetime
type request struct {
	Original       string     `json:"original"`
	Alias          string     `json:"alias"`
	CustomName     string     `json:"customName"`
	ClicksToExpire *int64     `json:"clicksToExpire,omitempty"`
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
}

func (r *request) ToModel(userId int64) *model.Link {
	return &model.Link{
		CreatedBy:      userId,
		Original:       r.Original,
		Alias:          r.Alias,
		CustomName:     r.CustomName,
		ExpirationDate: r.ExpirationDate,
		ClicksToExpire: r.ClicksToExpire,
	}
}

func (h *Handler) validateRequest(r *request) error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Original, validation.Required, validation.Length(1, h.OriginalMaxLen), is.URL),
		validation.Field(&r.Alias, validation.Length(1, h.AliasMaxLen)),
		validation.Field(&r.CustomName, validation.Length(1, h.CustomNameMaxLen)),
		validation.Field(&r.ClicksToExpire, validation.By(valkit.IsPositive())),
		validation.Field(&r.ExpirationDate, validation.By(valkit.IsFutureDate())),
	)
}
