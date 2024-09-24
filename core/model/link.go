package model

import (
	"time"
)

var (
	UnlimitedClicks int64 = -1
	AnonUser        int64 = -1
	NoExpireDate          = time.Time{}
	NeverVisited          = time.Time{}
)

type Link struct {
	Id                 int64
	CreatedBy          int64
	Original           string
	Alias              string
	CustomName         string
	ClicksCount        int64
	LastAccessTime     time.Time
	ClicksToExpiration int64
	ExpirationDate     time.Time
	Archived           bool
	CreatedAt          time.Time
}

//func (l *Link) CreatedByAnon() bool {
//	return l.CreatedBy == AnonUser
//}
//
//func (l *Link) IsExpired() bool {
//	if l.Archived {
//		fa
//	}
//
//	if l.ExpirationDate != NoExpireDate && time.Until(l.ExpirationDate) <= 0 {
//		return true
//	}
//
//	if l.ClicksToExpiration != UnlimitedClicks && l.ClicksToExpiration <= l.ClicksCount {
//		return true
//	}
//
//	return false
//}

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
