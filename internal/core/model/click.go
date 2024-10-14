package model

import (
	"net"
	"time"
)

type AdType string

const (
	AdTypeVideo = "video"
	AdTypeFile  = "file"
)

type ClickStatus string

const (
	ClickStatusOpened    = "opened"
	ClickStatusCompleted = "completed"
)

type ClickMetadata struct {
	IP         net.IP
	UserAgent  string
	AccessTime time.Time
}

type Click struct {
	Id         int64
	LinkId     int64
	AdSourceId int64
	AdType     AdType
	Status     ClickStatus
	Metadata   ClickMetadata
}

type GetClicksParams struct {
	Pagination Pagination
	Sort       Sort
	UserId     int64
	LinkId     int64
}
