package entity

import (
	"shortener/internal/core/model"
	"time"
)

type Click struct {
	Id         int64     `db:"id"`
	LinkId     int64     `db:"link_id"`
	IP         string    `db:"ip"`
	AccessTime time.Time `db:"access_time"`
	Status     int8      `db:"ad_status"`
}

func (c *Click) ToModel() *model.Click {
	return &model.Click{
		Metadata: &model.ClickMetadata{
			IP:         c.IP,
			AccessTime: c.AccessTime,
		},
		Id:     c.Id,
		LinkId: c.LinkId,
		Status: model.AdStatus(c.Status),
	}
}

func ClickFromModel(c *model.Click) *Click {
	return &Click{
		Id:         c.Id,
		LinkId:     c.LinkId,
		IP:         c.Metadata.IP,
		AccessTime: c.Metadata.AccessTime,
		Status:     int8(c.Status),
	}
}
