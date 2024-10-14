package response

import (
	"shortener/internal/core/model"
	"time"
)

type Click struct {
	Id         int64     `json:"id"`
	UserAgent  string    `json:"user_agent"`
	IP         string    `json:"ip"`
	Status     string    `json:"status"`
	AdType     string    `json:"ad_type"`
	AccessTime time.Time `json:"access_time"`
}

func ClickFromModel(click *model.Click) Click {
	return Click{
		Id:         click.Id,
		UserAgent:  click.Metadata.UserAgent,
		IP:         click.Metadata.IP.String(),
		AccessTime: click.Metadata.AccessTime,
		Status:     string(click.Status),
		AdType:     string(click.AdType),
	}
}
