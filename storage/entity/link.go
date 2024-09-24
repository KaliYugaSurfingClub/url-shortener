package entity

import (
	"database/sql"
	"time"
	"url_shortener/core/model"
)

type Link struct {
	Id                 int64         `db:"id"`
	CreatedBy          sql.NullInt64 `db:"created_by"`
	Original           string        `db:"original"`
	Alias              string        `db:"alias"`
	CustomName         string        `db:"custom_name"`
	ClicksCount        int64         `db:"clicks_count"`
	LastAccessTime     sql.NullTime  `db:"last_access_time"`
	ExpirationDate     sql.NullTime  `db:"expiration_date"`
	ClicksToExpiration sql.NullInt64 `db:"clicks_to_expiration"`
	Archived           bool          `db:"archived"`
	CreatedAt          time.Time     `db:"created_at"`
}

func (l *Link) ToModel() *model.Link {
	clicksToExpiration := model.UnlimitedClicks
	expirationDate := model.NoExpireDate
	createdBy := model.AnonUser
	lastAccessTime := model.NeverVisited

	if l.ClicksToExpiration.Valid {
		clicksToExpiration = l.ClicksToExpiration.Int64
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
		Id:                 l.Id,
		CreatedBy:          createdBy,
		Original:           l.Original,
		Alias:              l.Alias,
		CustomName:         l.CustomName,
		ClicksCount:        l.ClicksCount,
		LastAccessTime:     lastAccessTime,
		ClicksToExpiration: clicksToExpiration,
		ExpirationDate:     expirationDate,
		Archived:           l.Archived,
		CreatedAt:          l.CreatedAt,
	}
}

//todo clousure

func CreatedByToSql(id int64) sql.NullInt64 {
	if id == model.AnonUser {
		return sql.NullInt64{Valid: false}
	}

	return sql.NullInt64{Valid: true, Int64: id}
}

func ClicksToExpirationToSql(clicks int64) sql.NullInt64 {
	if clicks == model.UnlimitedClicks {
		return sql.NullInt64{Valid: false}
	}

	return sql.NullInt64{Valid: true, Int64: clicks}
}

func ExpirationDateToSql(date time.Time) sql.NullTime {
	if date == model.NoExpireDate {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Valid: true, Time: date}
}
