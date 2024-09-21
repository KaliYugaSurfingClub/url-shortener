package model

import (
	"time"
)

//todo add ORM or something like ORM to migrations

var (
	UnlimitedClicks int64 = -1
	NoExpireDate          = time.Time{}
)

type Link struct {
	Id          int64
	CreatedBy   int64
	Original    string
	Alias       string
	ClicksCount int64
	LastAccess  time.Time
	ExpireDate  time.Time
	MaxClicks   int64
}

type Click struct {
	Id         int64
	LinkId     int64
	AccessTime time.Time
	FullAD     bool
	IP         string
}

type User struct {
	id           int64
	email        string
	username     string
	passwordHash string
}
