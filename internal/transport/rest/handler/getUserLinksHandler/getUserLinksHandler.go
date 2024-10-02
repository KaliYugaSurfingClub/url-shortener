package getUserLinksHandler

import (
	"context"
	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/response"
	"shortener/internal/utils/valkit"
)

type request struct {
	Type        string `json:"type,omitempty"`
	Constraints string `json:"constraints,omitempty"`
	SortBy      string `json:"sortBy"`
	Order       string `json:"order"`
	Page        int64  `json:"page"`
	Size        int64  `json:"size"`
}

type data struct {
	TotalCount int64           `json:"totalCount"`
	Links      []response.Link `json:"links"`
}

type provider interface {
	GetUsersLinks(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, int64, error)
}

func New(provider provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

		userId, _ := mw.ExtractUserID(r.Context())

		params, err := paramsFromRequest(r)
		if err != nil {
			log.Error("invalid request", mw.ErrAttr(err))
			render.JSON(w, r, response.NewError(err))
			return
		}

		links, totalCount, err := provider.GetUsersLinks(r.Context(), userId, *params)
		if err != nil {
			log.Error("cannot get user links", mw.ErrAttr(err))
			render.JSON(w, r, response.NewInternalError())
			return
		}

		render.JSON(w, r, response.NewOk(data{
			TotalCount: totalCount,
			Links:      funk.Map(links, response.LinkFromModel).([]response.Link),
		}))
	}
}

func paramsFromRequest(r *http.Request) (*model.GetLinksParams, error) {
	defer r.Body.Close()

	var req request

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		return nil, err
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	return &model.GetLinksParams{
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
	}, nil
}

func (r *request) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Type, validation.By(valkit.ContainsInMap(types))),
		validation.Field(&r.Constraints, validation.By(valkit.ContainsInMap(constraints))),
		validation.Field(&r.SortBy, validation.By(valkit.ContainsInMap(sortBy))),
		validation.Field(&r.Order, validation.By(valkit.ContainsInMap(order))),
		validation.Field(&r.Page, validation.By(valkit.IsPositive())),
		validation.Field(&r.Size, validation.By(valkit.IsPositive())),
	)
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
	"createdAt":       model.SortByCreatedAt,
	"customName":      model.SortByCustomName,
	"clicksCount":     model.SortByClicksCount,
	"lastAccess":      model.SortByLastAccess,
	"expirationDate":  model.SortByExpirationDate,
	"leftClicksCount": model.SortByLeftClicksCount,
}

var order = map[string]model.Order{
	"asc":  model.Asc,
	"desc": model.Desc,
}
