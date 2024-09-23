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
	Id             int64
	CreatedBy      int64
	ClicksCount    int64
	MaxClicks      int64
	Original       string
	Alias          string
	LastAccessTime time.Time
	ExpirationDate time.Time
	CreatedAt      time.Time
}

func (l *Link) CreatedByAnon() bool {
	return l.CreatedBy == AnonUser
}

func (l *Link) IsExpired() bool {
	if l.ExpirationDate != NoExpireDate && time.Until(l.ExpirationDate) <= 0 {
		return true
	}

	if l.MaxClicks != UnlimitedClicks && l.MaxClicks <= l.ClicksCount {
		return true
	}

	return false
}

type Order int8

const (
	Asc Order = iota
	Desc
)

type TypeLink int8

const (
	TypeAny TypeLink = iota
	TypeActual
	TypeExpired
)

type ConstraintLink int8

const (
	ConstraintAny ConstraintLink = iota //means withMax withDate withoutAnything
	ConstraintMaxClicks
	ConstraintExpirationDate
	ConstraintWith
	ConstraintWithout //without any constraint
)

type ColumnLink int8

const (
	ColumnCreatedAt ColumnLink = iota
	ColumnAlias
	ColumnClicksCount
	ColumnLastAccess
	ColumnTimeToExpire
	ColumnClicksToExpire
)

type GetLinksParams struct {
	Type        TypeLink
	Constraints ConstraintLink
	Column      ColumnLink
	Order       Order
}
