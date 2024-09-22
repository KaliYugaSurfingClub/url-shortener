package entity

import (
	"database/sql"
	"time"
	"url_shortener/core/model"
)

type User struct {
	Id           int64  `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

type Link struct {
	Id             int64         `db:"id"`
	CreatedBy      int64         `db:"created_by"`
	Original       string        `db:"original"`
	Alias          string        `db:"alias"`
	ClicksCount    int64         `db:"clicks_count"`
	LastAccessTime sql.NullTime  `db:"last_access_time"`
	ExpirationDate sql.NullTime  `db:"expiration_date"`
	MaxClicks      sql.NullInt64 `db:"max_clicks"`
}

type Click struct {
	Id         int64     `db:"id"`
	LinkId     int64     `db:"link_id"`
	IP         string    `db:"ip"`
	FullAd     bool      `db:"full_ad"`
	AccessTime time.Time `db:"access_time"`
}

func (l *Link) ToModel() *model.Link {
	maxClicks := model.UnlimitedClicks
	expirationDate := model.NoExpireDate

	if l.MaxClicks.Valid {
		maxClicks = l.MaxClicks.Int64
	}

	if l.ExpirationDate.Valid {
		expirationDate = l.ExpirationDate.Time
	}

	return &model.Link{
		Id:             l.Id,
		CreatedBy:      l.CreatedBy,
		Original:       l.Original,
		Alias:          l.Alias,
		ClicksCount:    l.ClicksCount,
		MaxClicks:      maxClicks,
		ExpirationDate: expirationDate,
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
