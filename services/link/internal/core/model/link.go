package model

import (
	"time"
)

type Link struct {
	Id             int64
	CreatedBy      int64
	Original       string
	Alias          string
	CustomName     string
	ClicksCount    int64
	Archived       bool
	LastAccessTime *time.Time
	CreatedAt      time.Time
}

type GetLinksParams struct {
	Sort       Sort
	Pagination Pagination
	UserId     int64
	Archived   bool
}
