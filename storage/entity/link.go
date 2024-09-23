package entity

import (
	"database/sql"
	"time"
	"url_shortener/core/model"
)

type Link struct {
	Id             int64         `db:"id"`
	CreatedBy      sql.NullInt64 `db:"created_by"`
	Original       string        `db:"original"`
	Alias          string        `db:"alias"`
	ClicksCount    int64         `db:"clicks_count"`
	LastAccessTime sql.NullTime  `db:"last_access_time"`
	ExpirationDate sql.NullTime  `db:"expiration_date"`
	MaxClicks      sql.NullInt64 `db:"max_clicks"`
	CreatedAt      time.Time     `db:"created_at"`
}

func (l *Link) ToModel() *model.Link {
	maxClicks := model.UnlimitedClicks
	expirationDate := model.NoExpireDate
	createdBy := model.AnonUser
	lastAccessTime := model.NeverVisited

	if l.MaxClicks.Valid {
		maxClicks = l.MaxClicks.Int64
	}

	if l.ExpirationDate.Valid {
		expirationDate = l.ExpirationDate.Time
	}

	if l.CreatedBy.Valid {
		createdBy = l.CreatedBy.Int64
	}

	if l.LastAccessTime.Valid {
		lastAccessTime = l.LastAccessTime.Time
	}

	return &model.Link{
		Id:             l.Id,
		CreatedBy:      createdBy,
		Original:       l.Original,
		Alias:          l.Alias,
		ClicksCount:    l.ClicksCount,
		LastAccessTime: lastAccessTime,
		MaxClicks:      maxClicks,
		ExpirationDate: expirationDate,
		CreatedAt:      l.CreatedAt,
	}
}

func SqlMaxClicks(clicks int64) sql.NullInt64 {
	if clicks == model.UnlimitedClicks {
		return sql.NullInt64{Valid: false}
	}

	return sql.NullInt64{Valid: true, Int64: clicks}
}

func SqlExpirationDate(date time.Time) sql.NullTime {
	if date == model.NoExpireDate {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Valid: true, Time: date}
}
