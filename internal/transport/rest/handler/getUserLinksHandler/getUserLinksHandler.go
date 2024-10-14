package getUserLinksHandler

import (
	"context"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/request"
	"shortener/internal/transport/rest/response"
	"shortener/internal/utils/valkit"
)

type provider interface {
	GetUserLinks(ctx context.Context, params model.GetLinksParams) ([]*model.Link, int64, error)
}

type Handler struct {
	provider provider
}

type data struct {
	TotalCount int64           `json:"TotalCount"`
	Links      []response.Link `json:"links"`
}

func New(provider provider) *Handler {
	return &Handler{
		provider: provider,
	}
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

	urlParams := &UrlParams{}
	if err := request.DecodeURLParams(urlParams, r.URL.Query()); err != nil {
		log.Error("unable to decode URL urlParams", mw.ErrAttr(err))
		render.JSON(w, r, response.WithInternalError())
		return
	}

	if err := urlParams.Validate(); err != nil {
		log.Error("invalid url urlParams", mw.ErrAttr(err))
		render.JSON(w, r, response.WithValidationErrors(err))
		return
	}

	params := urlParams.ToModel()
	params.UserId, _ = mw.ExtractUserID(r.Context())

	links, totalCount, err := h.provider.GetUserLinks(r.Context(), params)
	if err != nil {
		log.Error("cannot get user links", mw.ErrAttr(err))
		render.JSON(w, r, response.WithInternalError())
		return
	}

	render.JSON(w, r, response.WithOk(data{
		TotalCount: totalCount,
		Links:      funk.Map(links, response.LinkFromModel).([]response.Link),
	}))
}

type UrlParams struct {
	Archived string `schema:"archived" json:"archived"`
	request.Pagination
	request.Sort
}

func (p *UrlParams) Validate() error {
	rules := []*validation.FieldRules{
		validation.Field(&p.Archived, validation.By(valkit.ContainsInMap(request.BoolMap))),
	}

	return request.Validate(p, rules, p.SortRules(sortBy), p.PaginationRules())
}

func (p *UrlParams) ToModel() model.GetLinksParams {
	return model.GetLinksParams{
		Archived:   request.BoolMap[p.Archived],
		Sort:       p.SortToModel(sortBy),
		Pagination: p.PaginationToModel(),
	}
}

var sortBy = map[string]model.SortBy{
	"created_at":   model.SortByCreatedAt,
	"custom_name":  model.SortLinksByCustomName,
	"clicks_count": model.SortLinksByClicksCount,
	"last_access":  model.SortLinksByLastAccess,
}
