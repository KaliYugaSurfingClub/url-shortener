package entity

import (
	"database/sql"
	"time"
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
