package linkRepo

import (
	"shortener/internal/core/model"
	"shortener/internal/storage/entity"
	"slices"
	"strings"
)

var (
	isArchivedLinkSql = ` archived = TRUE `
	isExpiredLinkSql  = `
		((expiration_date IS NOT NULL AND expiration_date <= CURRENT_TIMESTAMP) OR
		(clicks_to_expiration IS NOT NULL AND clicks_count >= clicks_to_expiration)) AND
		archived = FALSE 
	`
	isActiveLinkSql = `
		(expiration_date IS NULL OR expiration_date > CURRENT_TIMESTAMP) AND
		(clicks_to_expiration IS NULL OR clicks_count < clicks_to_expiration) AND
		archived = FALSE
	`
	isInactiveLinkSql = `
		(expiration_date IS NULL OR expiration_date > CURRENT_TIMESTAMP) OR
		(clicks_to_expiration IS NULL OR clicks_count < clicks_to_expiration) OR
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
		model.ConstraintWith:    "(clicks_to_expiration IS NOT NULL OR expiration_date IS NOT NULL)",
		model.ConstraintWithout: "clicks_count IS NULL AND expiration_date IS NULL",
		model.ConstraintClicks:  "clicks_to_expiration IS NOT NULL",
		model.ConstraintDate:    "expiration_date IS NOT NULL",
	}

	columnSql = map[model.SortByLink]string{
		model.SortByCreatedAt:       " created_at ",
		model.SortByCustomName:      " custom_name ",
		model.SortByClicksCount:     " clicks_count ",
		model.SortByLastAccess:      " last_access_time ",
		model.SortByLeftClicksCount: " COALESCE(clicks_to_expiration - clicks_count, -1) ",
		model.SortByExpirationDate:  " expiration_date ",
	}
)

//todo add pagination

func withGetParams(baseQuery string, params model.GetLinksParams) string {
	var query strings.Builder
	query.WriteString(baseQuery)

	conditions := make([]string, 0)

	conditions = append(conditions, typeSql[params.Type])
	conditions = append(conditions, constrainsSql[params.Constraints])

	conditions = slices.DeleteFunc(conditions, func(c string) bool { return c == "" })

	if len(conditions) > 0 {
		query.WriteString(" AND ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" ORDER BY ")
	query.WriteString(columnSql[params.SortBy])
	query.WriteString(entity.OrderToStr(params.Order))

	return query.String()
}

func activeOnly(baseQuery string) string {
	return baseQuery + " AND " + isActiveLinkSql
}
