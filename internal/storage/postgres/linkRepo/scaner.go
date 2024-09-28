package linkRepo

import (
	"github.com/jackc/pgx/v5"
	"shortener/internal/core/model"
)

func linkFromRow(row pgx.Row) (*model.Link, error) {
	link := &model.Link{}

	err := row.Scan(
		&link.Id, &link.CreatedBy, &link.Original, &link.Alias, &link.CustomName,
		&link.ClicksCount, &link.LastAccessTime, &link.ExpirationDate,
		&link.ClicksToExpire, &link.Archived, &link.CreatedAt,
	)

	return link, err
}
