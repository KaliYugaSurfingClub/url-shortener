package response

import (
	"shortener/internal/core/model"
	"time"
)

type Click struct {
	Id         int64          `json:"id"`
	Status     model.AdStatus `json:"status"`
	UserAgent  string         `json:"user_agent"`
	AccessTime time.Time      `json:"access_time"`
}

func ClickFromModel(click *model.Click) Click {
	return Click{
		Id:         click.Id,
		Status:     click.Status,
		UserAgent:  click.Metadata.UserAgent,
		AccessTime: click.Metadata.AccessTime,
	}
}
