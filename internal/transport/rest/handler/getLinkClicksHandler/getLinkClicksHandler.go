package getLinkClicksHandler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
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

	log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

	params, err := h.paramsFromQuery(r.URL.Query())
	if err != nil {
		log.Error("invalid url params", mw.ErrAttr(err))
		render.JSON(w, r, response.WithError(err))
		return
	}

	linkId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	clicks, totalCount, err := h.provider.GetLinkClicks(r.Context(), linkId, params)
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

func (h *Handler) paramsFromQuery(query url.Values) (model.GetClicksParams, error) {
	var ok bool

	params := model.GetClicksParams{}

	if params.Order, ok = order[query.Get("order")]; !ok {
		return params, fmt.Errorf("invalid query type")
	}

	if params.Pagination.Page, ok = positiveIntFromUrl(query, "page", 1); !ok {
		return params, fmt.Errorf("invalid query type")
	}

	if params.Pagination.Size, ok = positiveIntFromUrl(query, "size", h.defaultPageSize); !ok {
		return params, fmt.Errorf("invalid query type")
	}

	return params, nil
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
