package getLinkClicksHandler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/request"
	"shortener/internal/transport/rest/response"
	"strconv"
)

type provider interface {
	GetLinkClicks(ctx context.Context, params model.GetClicksParams) ([]*model.Click, int64, error)
}

type data struct {
	TotalCount int64
	Clicks     []response.Click
}

func New(provider provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.GetLinkClicks")

		urlParams := &UrlParams{}
		if err := request.DecodeURLParams(urlParams, r.URL.Query()); err != nil {
			response.Error(w, log, err)
			return
		}

		if err := urlParams.Validate(); err != nil {
			response.Error(w, log, err)
			return
		}

		params := urlParams.ToModel()
		params.UserId, _ = mw.ExtractUserID(r.Context())
		params.LinkId, _ = strconv.ParseInt(chi.URLParam(r, "linkId"), 10, 64)

		clicks, totalCount, err := provider.GetLinkClicks(r.Context(), params)
		if err != nil {
			response.Error(w, log, err)
			return
		}

		response.Ok(w, data{
			TotalCount: totalCount,
			Clicks:     funk.Map(clicks, response.ClickFromModel).([]response.Click),
		})
	}
}

type UrlParams struct {
	request.Sort
	request.Pagination
}

func (p *UrlParams) Validate() error {
	return request.Validate(p, p.SortRules(sortBy), p.PaginationRules())
}

func (p *UrlParams) ToModel() model.GetClicksParams {
	return model.GetClicksParams{
		Sort:       p.SortToModel(sortBy),
		Pagination: p.PaginationToModel(),
	}
}

var sortBy = map[string]model.SortBy{
	"access_time": model.SortClickByAccessTime,
}
