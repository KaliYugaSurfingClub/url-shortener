package linkRepo

import (
	"slices"
	"strings"
	"url_shortener/core/model"
	"url_shortener/storage/sqlite"
)

var (
	isActualLinkSql = `
		(expiration_date IS NULL OR expiration_date > CURRENT_TIMESTAMP) AND
		(max_clicks IS NULL OR clicks_count < max_clicks)
	`

	isExpiredLinkSql = ` NOT ( ` + isActualLinkSql + ") "

	typeSql = map[model.TypeLink]string{
		model.TypeAny:     "",
		model.TypeActual:  isActualLinkSql,
		model.TypeExpired: isExpiredLinkSql,
	}

	constrainsSql = map[model.ConstraintLink]string{
		model.ConstraintAny:            "",
		model.ConstraintWith:           "max_clicks IS NOT NULL AND expiration_date IS NOT NULL",
		model.ConstraintWithout:        "clicks_count IS NULL AND expiration_date IS NULL",
		model.ConstraintMaxClicks:      "max_clicks IS NOT NULL",
		model.ConstraintExpirationDate: "expiration_date IS NOT NULL",
	}

	columnSql = map[model.ColumnLink]string{
		model.ColumnCreatedAt:      " created_at ",
		model.ColumnAlias:          " alias ",
		model.ColumnClicksCount:    " clicks_count ",
		model.ColumnLastAccess:     " last_access_time ",
		model.ColumnClicksToExpire: " max_clicks - clicks_count ",
		model.ColumnTimeToExpire:   " expiration_date - CURRENT_TIMESTAMP ",
	}
)

//add pagination

func withGetParams(baseQuery string, params model.GetLinksParams) string {
	var query strings.Builder
	query.WriteString(baseQuery)

	conditions := make([]string, 0)

	conditions = append(conditions, typeSql[params.Type])
	conditions = append(conditions, constrainsSql[model.ConstraintWith])

	conditions = slices.DeleteFunc(conditions, func(c string) bool { return c == "" })

	if len(conditions) > 0 {
		query.WriteString(" AND ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" ORDER BY ")
	query.WriteString(columnSql[params.Column])
	query.WriteString(sqlite.OrderToStr(params.Order))

	return query.String()
}

func actualOnly(baseQuery string) string {
	return baseQuery + " AND " + isActualLinkSql
}
