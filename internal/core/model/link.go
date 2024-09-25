package model

import (
	"time"
)

type Link struct {
	Id                 int64
	CreatedBy          int64
	Original           string
	Alias              string
	CustomName         string
	ClicksCount        int64
	LastAccessTime     *time.Time
	ExpirationDate     *time.Time
	ClicksToExpiration *int64
	Archived           bool
	CreatedAt          time.Time
}

type Order int8

const (
	Asc Order = iota
	Desc
)

type LinkType int8

const (
	TypeAny LinkType = iota
	TypeActive
	TypeInactive
	TypeExpired
	TypeArchived
)

type SortByLink int8

const (
	SortByCreatedAt SortByLink = iota
	SortByCustomName
	SortByClicksCount
	SortByLastAccess
	SortByExpirationDate
	SortByLeftClicksCount
)

type LinkConstraints int8

const (
	ConstraintAny LinkConstraints = iota
	ConstraintClicks
	ConstraintDate
	ConstraintWith
	ConstraintWithout
)

type GetLinksParams struct {
	Type        LinkType
	Constraints LinkConstraints
	SortBy      SortByLink
	Order       Order
}
