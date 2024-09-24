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
	res := &model.Link{
		Id:                 l.Id,
		Original:           l.Original,
		Alias:              l.Alias,
		CustomName:         l.CustomName,
		ClicksCount:        l.ClicksCount,
		Archived:           l.Archived,
		CreatedAt:          l.CreatedAt,
		CreatedBy:          toNullableInt64(l.CreatedBy),
		LastAccessTime:     toNullableTime(l.LastAccessTime),
		ExpirationDate:     toNullableTime(l.ExpirationDate),
		ClicksToExpiration: toNullableInt64(l.ClicksToExpiration),
	}

	if res.CustomName == "" {
		res.CustomName = res.Alias
	}

	return res
}

func ModelToLink(m *model.Link) *Link {
	return &Link{
		Id:                 m.Id,
		Original:           m.Original,
		Alias:              m.Alias,
		CustomName:         m.CustomName,
		ClicksCount:        m.ClicksCount,
		Archived:           m.Archived,
		CreatedAt:          m.CreatedAt,
		CreatedBy:          fromNullableInt64(m.CreatedBy),
		LastAccessTime:     fromNullableTime(m.LastAccessTime),
		ExpirationDate:     fromNullableTime(m.ExpirationDate),
		ClicksToExpiration: fromNullableInt64(m.ClicksToExpiration),
	}
}

func fromNullableInt64(ptr *int64) sql.NullInt64 {
	if ptr != nil {
		return sql.NullInt64{Int64: *ptr, Valid: true}
	}
	return sql.NullInt64{Valid: false}
}

func fromNullableString(ptr *string) sql.NullString {
	if ptr != nil {
		return sql.NullString{String: *ptr, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func fromNullableTime(ptr *time.Time) sql.NullTime {
	if ptr != nil {
		return sql.NullTime{Time: *ptr, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

func toNullableInt64(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

func toNullableTime(n sql.NullTime) *time.Time {
	if n.Valid {
		return &n.Time
	}
	return nil
}