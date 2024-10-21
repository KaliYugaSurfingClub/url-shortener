package getUserLinksHandler

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/request"
	"shortener/internal/transport/rest/response"
	"shortener/internal/utils/valkit"
)

type LinksProvider interface {
	GetUserLinks(ctx context.Context, params model.GetLinksParams) ([]*model.Link, int64, error)
}

type data struct {
	TotalCount int64           `json:"TotalCount"`
	Links      []response.Link `json:"links"`
}

func New(provider LinksProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

		urlParams := &UrlParams{}
		if err := request.DecodeURLParams(urlParams, r.URL.Query()); err != nil {
			//todo decode error
			return
		}

		if err := urlParams.Validate(); err != nil {
			response.Error(w, log, err)
			return
		}

		params := urlParams.ToModel()
		params.UserId, _ = mw.ExtractUserID(r.Context())

		links, totalCount, err := provider.GetUserLinks(r.Context(), params)
		if err != nil {
			response.Error(w, log, err)
			return
		}

		response.Ok(w, data{
			TotalCount: totalCount,
			Links:      funk.Map(links, response.LinkFromModel).([]response.Link),
		})
	}
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
