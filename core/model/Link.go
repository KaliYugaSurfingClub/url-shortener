package model

import "time"

var (
	UnlimitedClicks int64 = -1
	NoExpireDate          = time.Time{}
)

type Link struct {
	Id             int64
	CreatedBy      int64
	Original       string
	Alias          string
	ClicksCount    int64
	LastAccess     time.Time
	ExpirationDate time.Time
	MaxClicks      int64
}

func (l *Link) IsExpired() bool {
	validTime := l.ExpirationDate == NoExpireDate || time.Until(l.ExpirationDate) > 0
	validClicks := l.MaxClicks == UnlimitedClicks && l.MaxClicks >= l.ClicksCount

	return !validTime || validClicks
}
