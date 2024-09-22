package entity

import (
	"time"
)

type User struct {
	Id           int64  `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

type Click struct {
	Id         int64     `db:"id"`
	LinkId     int64     `db:"link_id"`
	IP         string    `db:"ip"`
	FullAd     bool      `db:"full_ad"`
	AccessTime time.Time `db:"access_time"`
}
