package model

import "time"

type AdStatus int8

const (
	AdStarted = iota
	AdClosed
	AdCompleted
)

type ClickMetadata struct {
	IP         string
	UserAgent  string
	AccessTime time.Time
}

type Click struct {
	Id       int64
	LinkId   int64
	Status   AdStatus
	Metadata *ClickMetadata
}
