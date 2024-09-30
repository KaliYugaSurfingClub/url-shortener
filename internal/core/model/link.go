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
	LastAccessTime *time.Time
	ExpirationDate *time.Time
	ClicksToExpire *int64
	Archived       bool
	CreatedAt      time.Time
}

type Order int8

const (
	Asc Order = iota
	Desc
)

type LinkSortBy int8

const (
	SortByCreatedAt LinkSortBy = iota
	SortByCustomName
	SortByClicksCount
	SortByLastAccess
	SortByExpirationDate
	SortByLeftClicksCount
)

type LinkType int8

const (
	TypeAny LinkType = iota
	TypeActive
	TypeInactive
	TypeExpired
	TypeArchived
)

type LinkConstraints int8

const (
	ConstraintAny LinkConstraints = iota
	ConstraintClicks
	ConstraintDate
	ConstraintWith
	ConstraintWithout
)

type Pagination struct {
	Page int
	Size int
}

type LinkFilter struct {
	Type        LinkType
	Constraints LinkConstraints
}

type LinkSort struct {
	SortBy LinkSortBy
	Order  Order
}

type GetLinksParams struct {
	Filter     LinkFilter
	Sort       LinkSort
	Pagination Pagination
}
