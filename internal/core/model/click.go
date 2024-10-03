package model

import (
	"net"
	"time"
)

type AdStatus int8

const (
	AdStarted = iota
	AdClosed
	AdCompleted
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
	Status     AdStatus
	Order      Order
}
