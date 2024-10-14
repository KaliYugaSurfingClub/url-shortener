package builder

import (
	"shortener/internal/core/model"
	"strconv"
	"strings"
)

type Builder interface {
	Paginate(model.Pagination) Builder
	String() string
}

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

func (b *BaseBuilder) Paginate(params model.Pagination) Builder {
	offset := (params.Page - 1) * params.Size
	limit := params.Size

	b.Query.WriteString(" LIMIT ")
	b.Query.WriteString(strconv.FormatInt(limit, 10))
	b.Query.WriteString(" OFFSET ")
	b.Query.WriteString(strconv.FormatInt(offset, 10))

	return b
}

func (b *BaseBuilder) Sort(columnNames map[model.SortBy]string, sort model.Sort) Builder {
	column, ok := columnNames[sort.By]
	if !ok {
		panic("column not found")
	}

	b.Query.WriteString(" ORDER BY ")
	b.Query.WriteString(column)
	b.Query.WriteString(postgresOrder[sort.Order])

	return b
}

var postgresOrder = map[model.Order]string{
	model.Desc: "DESC",
	model.Asc:  "ASC",
}
