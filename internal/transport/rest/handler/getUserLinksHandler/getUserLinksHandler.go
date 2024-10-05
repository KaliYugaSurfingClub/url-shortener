package getUserLinksHandler

import (
	"context"
	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/request"
	"shortener/internal/transport/rest/response"
	"shortener/internal/utils/valkit"
)

type provider interface {
	GetUserLinks(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, int64, error)
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
	userId, _ := mw.ExtractUserID(r.Context())

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

	links, totalCount, err := h.provider.GetUserLinks(r.Context(), userId, params.ToModel())
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
	Type        string `schema:"type" json:"type"`
	Constraints string `schema:"constraints" json:"constraints"`
	SortBy      string `schema:"sort_by" json:"sort_by"`
	request.Pagination
	request.OrderDirection
}

func (p *UrlParams) Validate() error {
	rules := []*validation.FieldRules{
		validation.Field(&p.Type, validation.By(valkit.ContainsInMap(types))),
		validation.Field(&p.Constraints, validation.By(valkit.ContainsInMap(constraints))),
		validation.Field(&p.SortBy, validation.By(valkit.ContainsInMap(sortBy))),
	}

	return request.Validate(p, rules, p.OrderRules(), p.PaginationRules())
}

func (p *UrlParams) ToModel() model.GetLinksParams {
	return model.GetLinksParams{
		Filter: model.LinkFilter{
			Type:        types[p.Type],
			Constraints: constraints[p.Constraints],
		},
		Sort: model.LinkSort{
			By:    sortBy[p.SortBy],
			Order: p.OrderToModel(),
		},
		Pagination: p.PaginationToModel(),
	}
}

var types = map[string]model.LinkType{
	"any":      model.TypeAny,
	"active":   model.TypeActive,
	"inactive": model.TypeInactive,
	"expired":  model.TypeExpired,
	"archived": model.TypeArchived,
}

var constraints = map[string]model.LinkConstraints{
	"any":     model.ConstraintAny,
	"clicks":  model.ConstraintClicks,
	"date":    model.ConstraintDate,
	"with":    model.ConstraintWith,
	"without": model.ConstraintWithout,
}

var sortBy = map[string]model.LinkSortBy{
	"created_at":       model.SortByCreatedAt,
	"custom_name":      model.SortByCustomName,
	"clicks_count":     model.SortByClicksCount,
	"last_access":      model.SortByLastAccess,
	"expiration_date":  model.SortByExpirationDate,
	"left_clicksCount": model.SortByLeftClicksCount,
}
