package model

import "time"

var (
	UnlimitedClicks int64 = -1
	AnonUser        int64 = -1
	NoExpireDate          = time.Time{}
	NeverVisited          = time.Time{}
)

type Link struct {
	Id             int64
	CreatedBy      int64
	Original       string
	Alias          string
	ClicksCount    int64
	LastAccessTime time.Time
	ExpirationDate time.Time
	MaxClicks      int64
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
