package rest

import (
	"encoding/json"
	"github.com/KaliYugaSurfingClub/pkg/errs"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/schema"
	"link-service/internal/core/model"
	"link-service/internal/utils/valkit"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strconv"
)

func DecodeJSON(dst any, r *http.Request) error {
	const op errs.Op = "transport.rest.request.DecodeJSON"

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return errs.E(op, err, errs.InvalidRequest)
	}

	return nil
}

func DecodeURLParams(dst any, query url.Values) error {
	const op errs.Op = "transport.rest.request.DecodeURLParams"

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(dst, query); err != nil {
		return errs.E(op, err, errs.InvalidRequest)
	}

	return nil
}

type Pagination struct {
	Page string `schema:"page" json:"page"` // parse url params to string for custom validation errors
	Size string `schema:"size" json:"size"`
}

func (p *Pagination) PaginationRules() []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&p.Page, validation.By(valkit.StringNumIn(1, math.MaxInt))),
		validation.Field(&p.Size, validation.By(valkit.StringNumIn(0, math.MaxInt))),
	}
}

func (p *Pagination) PaginationToModel() (res model.Pagination) {
	res.Page, _ = strconv.ParseInt(p.Page, 10, 64)
	res.Size, _ = strconv.ParseInt(p.Size, 10, 64)

	return res
}

type Sort struct {
	Order string `schema:"order" json:"order"`
	By    string `schema:"sort_by" json:"sort_by"`
}

func (s *Sort) SortRules(sortBy map[string]model.SortBy) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&s.Order, validation.By(valkit.ContainsInMap(OrderMap))),
		validation.Field(&s.By, validation.By(valkit.ContainsInMap(sortBy))),
	}
}

func (s *Sort) SortToModel(sortBy map[string]model.SortBy) (res model.Sort) {
	return model.Sort{
		Order: OrderMap[s.Order],
		By:    sortBy[s.By],
	}
}

func Validate(ptr any, rules ...[]*validation.FieldRules) error {
	const op errs.Op = "transport.rest.request.Validate"

	err := validation.ValidateStruct(ptr, slices.Concat(rules...)...)
	if err == nil {
		return nil
	}

	return errs.E(op, err, errs.Validation)
}

var OrderMap = map[string]model.Order{
	"asc":  model.Asc,
	"desc": model.Desc,
}

var BoolMap = map[string]bool{
	"true":  true,
	"false": false,
}
