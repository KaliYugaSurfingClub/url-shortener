package getUserLinksHandler

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/thoas/go-funk"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/jsonValidator"
	"shortener/internal/transport/rest/mw"
	"shortener/internal/transport/rest/response"
	"strings"
)

type request struct {
	Type        string `json:"type" validate:"validate_type"`
	Constraints string `json:"constraints" validate:"validate_constraints"`
	SortBy      string `json:"sortBy" validate:"validate_sortBy"`
	Order       string `json:"order" validate:"validate_order"`
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
	valid := jsonValidator.New()

	valid.AddValidation(
		validationFuncFromMap("validate_type", types),
		validationFuncFromMap("validate_constraints", constraints),
		validationFuncFromMap("validate_sortBy", sortBy),
		validationFuncFromMap("validate_order", order),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

		userId, _ := mw.ExtractUserID(r.Context())

		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("cannot decode body", mw.ErrAttr(err))
			render.JSON(w, r, response.NewError(err))
			return
		}

		if err := valid.Validate(req); err != nil {
			log.Error("validation error", mw.ErrAttr(err))
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

func validationFuncFromMap[T any](name string, acceptable map[string]T) jsonValidator.ValidationFunc {
	fn := func(fl validator.FieldLevel) bool {
		_, ok := types[fl.Field().String()]
		return ok
	}

	keys := funk.Keys(acceptable).([]string)
	inBuckets := strings.Join(funk.Map(keys, func(s string) string { return s }).([]string), " ")

	return jsonValidator.ValidationFunc{
		Name: name,
		Fn:   fn,
		Err:  fmt.Errorf("should be one from %s", strings.TrimSpace(inBuckets)),
	}
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
	"Any":     model.ConstraintAny,
	"Clicks":  model.ConstraintClicks,
	"Date":    model.ConstraintDate,
	"With":    model.ConstraintWith,
	"Without": model.ConstraintWithout,
}

var sortBy = map[string]model.LinkSortBy{
	"CreatedAt":       model.SortByCreatedAt,
	"CustomName":      model.SortByCustomName,
	"ClicksCount":     model.SortByClicksCount,
	"LastAccess":      model.SortByLastAccess,
	"ExpirationDate":  model.SortByExpirationDate,
	"LeftClicksCount": model.SortByLeftClicksCount,
}

var order = map[string]model.Order{
	"asc":  model.Asc,
	"desc": model.Desc,
}

func paramsFromRequest(req *request) model.GetLinksParams {
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
