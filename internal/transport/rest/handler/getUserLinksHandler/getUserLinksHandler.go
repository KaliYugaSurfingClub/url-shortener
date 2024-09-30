package getUserLinksHandler

import (
	"context"
	"github.com/go-chi/render"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/response"
)

// todo all may be omit empty maybe
type request struct {
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
	SortBy      string `json:"sortBy"`
	Order       string `json:"order"`
	Page        int    `json:"page"`
	Size        int    `json:"size"`
}

type data struct {
	TotalCount int64           `json:"totalCount"`
	Links      []response.Link `json:"links"`
}

type LinksProvider interface {
	GetByUserId(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, error)
	GetCountByUserId(ctx context.Context, userId int64, params model.LinkFilter) (int64, error)
}

func New(provider LinksProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

		userId, _ := mw.ExtractUserID(r.Context())

		//todo validate and scan into model.Params in one func
		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("cannot decode body", mw.ErrAttr(err))
			render.JSON(w, r, response.NewError(err))
			return
		}

		params := paramsFromRequest(&req)

		totalCount, err := provider.GetCountByUserId(r.Context(), userId, params.Filter)
		if err != nil {
			log.Error("cannot get count of user links", mw.ErrAttr(err))
			render.JSON(w, r, response.NewError(err))
			return
		}

		links, err := provider.GetByUserId(r.Context(), userId, params)
		if err != nil {
			log.Error("cannot get user links", mw.ErrAttr(err))
			render.JSON(w, r, response.NewError(err))
			return
		}

		render.JSON(w, r, response.NewOk(data{
			TotalCount: totalCount,
			Links:      funk.Map(links, func(item *model.Link) response.Link { return response.LinkFromModel(item) }).([]response.Link),
		}))
	}
}

func paramsFromRequest(req *request) model.GetLinksParams {
	types := map[string]model.LinkType{
		"any":      model.TypeAny,
		"active":   model.TypeActive,
		"inactive": model.TypeInactive,
		"expired":  model.TypeExpired,
		"archived": model.TypeArchived,
	}

	constraints := map[string]model.LinkConstraints{
		"Any":     model.ConstraintAny,
		"Clicks":  model.ConstraintClicks,
		"Date":    model.ConstraintDate,
		"With":    model.ConstraintWith,
		"Without": model.ConstraintWithout,
	}

	sortBy := map[string]model.LinkSortBy{
		"CreatedAt":       model.SortByCreatedAt,
		"CustomName":      model.SortByCustomName,
		"ClicksCount":     model.SortByClicksCount,
		"LastAccess":      model.SortByLastAccess,
		"ExpirationDate":  model.SortByExpirationDate,
		"LeftClicksCount": model.SortByLeftClicksCount,
	}

	order := map[string]model.Order{
		"asc":  model.Asc,
		"desc": model.Desc,
	}

	return model.GetLinksParams{
		Filter: model.LinkFilter{
			Type:        types[req.Type],
			Constraints: constraints[req.Constraints],
		},
		Sort: model.LinkSort{
			SortBy: sortBy[req.SortBy],
			Order:  order[req.Order],
		},
		Pagination: model.Pagination{
			Page: req.Page,
			Size: req.Size,
		},
	}
}
