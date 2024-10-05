package request

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/schema"
	"math"
	"net/url"
	"shortener/internal/core/model"
	"shortener/internal/utils/valkit"
	"slices"
	"strconv"
)

func DecodeURLParams(dst any, query url.Values) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(dst, query); err != nil {
		return err
	}

	return nil
}

type Pagination struct {
	Page string `schema:"page" json:"page"`
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

// DO NOT RENAME!!!

type OrderDirection struct {
	Order string `schema:"order" json:"order"`
}

var OrderMap = map[string]model.Order{
	"asc":  model.Asc,
	"desc": model.Desc,
}

func (o *OrderDirection) OrderRules() []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&o.Order, validation.By(valkit.ContainsInMap(OrderMap))),
	}
}

func (o *OrderDirection) OrderToModel() (res model.Order) {
	return OrderMap[o.Order]
}

func Validate(ptr any, rules ...[]*validation.FieldRules) error {
	return validation.ValidateStruct(ptr, slices.Concat(rules...)...)
}
