package getUserLinksHandler

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/thoas/go-funk"
	"net/http"
	"net/url"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/response"
	"strconv"
)

type provider interface {
	GetUserLinks(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, int64, error)
}

type Handler struct {
	provider        provider
	defaultPageSize int64
}

type data struct {
	TotalCount int64           `json:"totalCount"`
	Links      []response.Link `json:"links"`
}

func New(provider provider, defaultPageSize int64) *Handler {
	return &Handler{
		provider:        provider,
		defaultPageSize: defaultPageSize,
	}
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

	userId, _ := mw.ExtractUserID(r.Context())

	params, err := h.paramsFromQuery(r.URL.Query())
	if err != nil {
		log.Error("invalid url params", mw.ErrAttr(err))
		render.JSON(w, r, response.WithError(err))
		return
	}

	links, totalCount, err := h.provider.GetUserLinks(r.Context(), userId, *params)
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

func (h *Handler) paramsFromQuery(query url.Values) (*model.GetLinksParams, error) {
	var ok bool

	params := &model.GetLinksParams{}

	if params.Filter.Type, ok = types[query.Get("type")]; !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	if params.Filter.Constraints, ok = constraints[query.Get("constraints")]; !ok {
		return nil, fmt.Errorf("invalid query constraints")
	}

	if params.Sort.By, ok = sortBy[query.Get("sort_by")]; !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	if params.Sort.Order, ok = order[query.Get("order")]; !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	if params.Pagination.Page, ok = positiveIntFromUrl(query, "page", 1); !ok {
		return nil, fmt.Errorf("invalid query type") //todo validation error
	}

	if params.Pagination.Size, ok = positiveIntFromUrl(query, "size", h.defaultPageSize); !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	return params, nil
}

var types = map[string]model.LinkType{
	"":         model.TypeAny,
	"any":      model.TypeAny,
	"active":   model.TypeActive,
	"inactive": model.TypeInactive,
	"expired":  model.TypeExpired,
	"archived": model.TypeArchived,
}

var constraints = map[string]model.LinkConstraints{
	"":        model.ConstraintAny,
	"any":     model.ConstraintAny,
	"clicks":  model.ConstraintClicks,
	"date":    model.ConstraintDate,
	"with":    model.ConstraintWith,
	"without": model.ConstraintWithout,
}

var sortBy = map[string]model.LinkSortBy{
	"":                 model.SortByCreatedAt,
	"created_at":       model.SortByCreatedAt,
	"custom_name":      model.SortByCustomName,
	"clicks_count":     model.SortByClicksCount,
	"last_access":      model.SortByLastAccess,
	"expiration_date":  model.SortByExpirationDate,
	"left_clicksCount": model.SortByLeftClicksCount,
}

var order = map[string]model.Order{
	"":     model.Asc,
	"asc":  model.Asc,
	"desc": model.Desc,
}

func positiveIntFromUrl(query url.Values, varName string, defaultValue int64) (int64, bool) {
	str := query.Get(varName)
	if str == "" {
		str = strconv.FormatInt(defaultValue, 10)
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}

	if num <= 0 {
		return 0, false
	}

	return int64(num), true
}
