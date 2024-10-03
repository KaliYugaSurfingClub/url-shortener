package clickRepo

import (
	"shortener/internal/core/model"
	"strconv"
	"strings"
)

type builder struct {
	query strings.Builder
}

func build(baseQuery string) *builder {
	res := new(builder)
	res.query.WriteString(baseQuery)

	return res
}

// todo
func (b *builder) Paginate(params model.Pagination) *builder {
	offset := (params.Page - 1) * params.Size
	limit := params.Size

	b.query.WriteString(" LIMIT ")
	b.query.WriteString(strconv.FormatInt(limit, 10))
	b.query.WriteString(" OFFSET ")
	b.query.WriteString(strconv.FormatInt(offset, 10))

	return b
}

func (b *builder) String() string {
	return b.query.String()
}
