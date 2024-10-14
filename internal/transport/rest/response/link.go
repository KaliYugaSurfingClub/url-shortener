package response

import (
	"shortener/internal/core/model"
	"time"
)

type Link struct {
	Id             int64      `json:"id"`
	CreatedBy      int64      `json:"createdBy"`
	Original       string     `json:"original"`
	Alias          string     `json:"alias"`
	CustomName     string     `json:"customName"`
	ClicksCount    int64      `json:"clicksCount"`
	LastAccessTime *time.Time `json:"lastAccessTime,omitempty"`
	Archived       bool       `json:"archived"`
	CreatedAt      time.Time  `json:"createdAt"`
}

func LinkFromModel(link *model.Link) Link {
	return Link{
		Id:             link.Id,
		CreatedBy:      link.CreatedBy,
		Original:       link.Original,
		Alias:          link.Alias,
		CustomName:     link.CustomName,
		ClicksCount:    link.ClicksCount,
		LastAccessTime: link.LastAccessTime,
		Archived:       link.Archived,
		CreatedAt:      link.CreatedAt,
	}
}
