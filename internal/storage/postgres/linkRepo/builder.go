package linkRepo

import (
	"shortener/internal/core/model"
	"shortener/internal/storage/postgres/builder"
	"slices"
	"strings"
)

type Builder struct {
	*builder.BaseBuilder
}

func newBuilder(baseQuery string) *Builder {
	return &Builder{builder.New(baseQuery)}
}

func (b *Builder) Filter(params model.LinkFilter) *Builder {
	conditions := make([]string, 0)

	conditions = append(conditions, typeSql[params.Type])
	conditions = append(conditions, constrainsSql[params.Constraints])

	conditions = slices.DeleteFunc(conditions, func(c string) bool { return c == "" })

	if len(conditions) > 0 {
		b.Query.WriteString(" AND ")
		b.Query.WriteString(strings.Join(conditions, " AND "))
	}

	return b
}

func (b *Builder) Sort(params model.LinkSort) *Builder {
	b.Query.WriteString(" ORDER BY ")
	b.Query.WriteString(sortBy[params.SortBy])
	b.Query.WriteString(OrderToStr(params.Order))

	return b
}

func (b *Builder) String() string {
	return b.Query.String()
}

func activeOnly(baseQuery string) string {
	return baseQuery + " AND " + isActiveLinkSql
}

func OrderToStr(order model.Order) string {
	if order == model.Desc {
		return "DESC"
	}

	return "ASC"
}

var (
	isArchivedLinkSql = ` archived = TRUE `
	isExpiredLinkSql  = `
		((expiration_date IS NOT NULL AND expiration_date <= CURRENT_TIMESTAMP) OR
		(clicks_to_expire IS NOT NULL AND clicks_count >= clicks_to_expire)) AND
		archived = FALSE 
	`
	isActiveLinkSql = `
		(expiration_date IS NULL OR expiration_date > CURRENT_TIMESTAMP) AND
		(clicks_to_expire IS NULL OR clicks_count < clicks_to_expire) AND
		archived = FALSE
	`
	isInactiveLinkSql = `
		(expiration_date IS NULL OR expiration_date > CURRENT_TIMESTAMP) OR
		(clicks_to_expire IS NULL OR clicks_count < clicks_to_expire) OR
		archived = TRUE 
	`

	typeSql = map[model.LinkType]string{
		model.TypeAny:      "",
		model.TypeActive:   isActiveLinkSql,
		model.TypeExpired:  isExpiredLinkSql,
		model.TypeArchived: isArchivedLinkSql,
		model.TypeInactive: isInactiveLinkSql,
	}

	constrainsSql = map[model.LinkConstraints]string{
		model.ConstraintAny:     "",
		model.ConstraintWith:    "(clicks_to_expire IS NOT NULL OR expiration_date IS NOT NULL)",
		model.ConstraintWithout: "clicks_count IS NULL AND expiration_date IS NULL",
		model.ConstraintClicks:  "clicks_to_expire IS NOT NULL",
		model.ConstraintDate:    "expiration_date IS NOT NULL",
	}

	sortBy = map[model.LinkSortBy]string{
		model.SortByCreatedAt:       " created_at ",
		model.SortByCustomName:      " custom_name ",
		model.SortByClicksCount:     " clicks_count ",
		model.SortByLastAccess:      " last_access_time ",
		model.SortByLeftClicksCount: " COALESCE(clicks_to_expire - clicks_count, -1) ",
		model.SortByExpirationDate:  " expiration_date ",
	}
)
