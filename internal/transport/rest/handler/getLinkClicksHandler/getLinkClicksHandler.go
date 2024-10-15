package getLinkClicksHandler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core"
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
			log.Error("unable to decode URL urlParams", mw.ErrAttr(err))
			render.JSON(w, r, response.WithInternalError())
			return
		}

		defer r.Body.Close()

		if err := urlParams.Validate(); err != nil {
			log.Error("invalid url urlParams", mw.ErrAttr(err))
			render.JSON(w, r, response.WithValidationErrors(err))
			return
		}

		params := urlParams.ToModel()
		//todo
		params.UserId, _ = mw.ExtractUserID(r.Context())
		params.LinkId, _ = strconv.ParseInt(chi.URLParam(r, "linkId"), 10, 64)

		clicks, totalCount, err := provider.GetLinkClicks(r.Context(), params)
		if errors.Is(err, core.ErrLinkNotFound) {
			log.Info("link not found", mw.ErrAttr(err))
			render.JSON(w, r, response.WithError(core.ErrLinkNotFound))
			return
		}
		if err != nil {
			log.Error("cannot get user links", mw.ErrAttr(err))
			render.JSON(w, r, response.WithInternalError())
			return
		}

		render.JSON(w, r, response.WithData(data{
			TotalCount: totalCount,
			Clicks:     funk.Map(clicks, response.ClickFromModel).([]response.Click),
		}))
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
