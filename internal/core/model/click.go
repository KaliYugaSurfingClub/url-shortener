package model

import (
	"net"
	"time"
)

type AdStatus string

const (
	AdStarted   = "started"
	AdClosed    = "closed"
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
}
