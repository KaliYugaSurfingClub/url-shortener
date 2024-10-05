package getLinkClicksHandler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/request"
	"shortener/internal/transport/rest/response"
	"strconv"
)

type provider interface {
	GetLinkClicks(ctx context.Context, linkId int64, params model.GetClicksParams) ([]*model.Click, int64, error)
}

type Handler struct {
	provider        provider
	defaultPageSize int64
}

type data struct {
	TotalCount int64
	Clicks     []response.Click
}

func New(provider provider, defaultPageSize int64) *Handler {
	return &Handler{
		provider:        provider,
		defaultPageSize: defaultPageSize,
	}
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log := mw.ExtractLog(r.Context(), "transport.rest.GetLinkClicks")

	params := &UrlParams{}
	if err := request.DecodeURLParams(params, r.URL.Query()); err != nil {
		log.Error("unable to decode URL params", mw.ErrAttr(err))
		render.JSON(w, r, response.WithInternalError())
		return
	}

	if err := params.Validate(); err != nil {
		log.Error("invalid url params", mw.ErrAttr(err))
		render.JSON(w, r, response.WithValidationErrors(err))
		return
	}

	linkId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	clicks, totalCount, err := h.provider.GetLinkClicks(r.Context(), linkId, params.ToModel())
	if err != nil {
		log.Error("cannot get user links", mw.ErrAttr(err))
		render.JSON(w, r, response.WithInternalError())
		return
	}

	render.JSON(w, r, response.WithOk(data{
		TotalCount: totalCount,
		Clicks:     funk.Map(clicks, response.ClickFromModel).([]response.Click),
	}))
}

type UrlParams struct {
	request.OrderDirection
	request.Pagination
}

func (p *UrlParams) Validate() error {
	return request.Validate(p, p.OrderRules(), p.PaginationRules())
}

func (p *UrlParams) ToModel() model.GetClicksParams {
	return model.GetClicksParams{
		Order:      p.OrderToModel(),
		Pagination: p.PaginationToModel(),
	}
}
