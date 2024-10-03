package clickRepo

import (
	"github.com/jackc/pgx/v5"
	"shortener/internal/core/model"
)

func clickFromRow(row pgx.Row) (*model.Click, error) {
	click := &model.Click{}

	err := row.Scan(
		&click.Id, &click.LinkId,
		&click.Metadata.UserAgent, &click.Metadata.IP, &click.Metadata.AccessTime,
		&click.Status,
	)

	return click, err
}
