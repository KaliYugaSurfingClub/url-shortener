package builder

import (
	"shortener/internal/core/model"
	"strconv"
	"strings"
)

type BaseBuilder struct {
	Query strings.Builder
}

func New(query string) *BaseBuilder {
	res := new(BaseBuilder)
	res.Query.WriteString(query)

	return res
}

func (b *BaseBuilder) String() string {
	return b.Query.String()
}

// todo
func (b *BaseBuilder) Paginate(params model.Pagination) *BaseBuilder {
	offset := (params.Page - 1) * params.Size
	limit := params.Size

	b.Query.WriteString(" LIMIT ")
	b.Query.WriteString(strconv.FormatInt(limit, 10))
	b.Query.WriteString(" OFFSET ")
	b.Query.WriteString(strconv.FormatInt(offset, 10))

	return b
}