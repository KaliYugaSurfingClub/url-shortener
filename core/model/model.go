package model

import (
	"time"
)

type Click struct {
	Id         int64
	LinkId     int64
	FullAD     bool
	IP         string
	AccessTime time.Time
}

type User struct {
	id           int64
	email        string
	username     string
	passwordHash string
}
