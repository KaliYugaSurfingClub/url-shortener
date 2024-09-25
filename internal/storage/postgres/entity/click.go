package entity

import (
	"shortener/internal/core/model"
	"shortener/internal/storage/postgres"
	"time"
)

type Click struct {
	Id         int64             `db:"id"`
	LinkId     int64             `db:"link_id"`
	Status     int8              `db:"ad_status"`
	UserAgent  string            `db:"user_agent"`
	IP         postgres.NullInet `db:"ip"`
	AccessTime time.Time         `db:"access_time"`
}

func (c *Click) ToModel() *model.Click {
	res := &model.Click{
		Metadata: model.ClickMetadata{
			AccessTime: c.AccessTime,
			UserAgent:  c.UserAgent,
		},
		Id:     c.Id,
		LinkId: c.LinkId,
		Status: model.AdStatus(c.Status),
	}

	if c.IP.Valid {
		res.Metadata.IP = c.IP.IP
	}

	return res
}

func ClickFromModel(c *model.Click) *Click {
	res := &Click{
		Id:         c.Id,
		LinkId:     c.LinkId,
		AccessTime: c.Metadata.AccessTime,
		UserAgent:  c.Metadata.UserAgent,
		Status:     int8(c.Status),
	}

	//todo add validation
	if c.Metadata.IP != nil {
		res.IP = postgres.NullInet{Valid: true, IP: c.Metadata.IP}
	}

	return res
}
