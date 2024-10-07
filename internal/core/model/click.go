package model

import (
	"net"
	"time"
)

type AdStatus string

// todo parse from db
const (
	AdStarted   = "started"
	AdClosed    = "closed"
	AdWatched   = "watched"
	AdCompleted = "completed"
)

type ClickMetadata struct {
	IP         net.IP
	UserAgent  string
	AccessTime time.Time
}

type Click struct {
	Id       int64
	LinkId   int64
	Status   AdStatus
	Metadata ClickMetadata
}

type GetClicksParams struct {
	Pagination Pagination
	Order      Order
	UserId     int64
	LinkId     int64
}
