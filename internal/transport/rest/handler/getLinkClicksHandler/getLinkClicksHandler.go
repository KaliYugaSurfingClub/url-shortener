package getLinkClicksHandler

import (
	"context"
	"shortener/internal/core/model"
)

type provider interface {
	GetLinkClicks(ctx context.Context, linkId int64, params model.GetClicksParams) ([]*model.Link, int64, error)
}

type Handler struct {
	provider        provider
	defaultPageSize int64
}

func New(provider provider, defaultPageSize int64) *Handler {
	return &Handler{
		provider:        provider,
		defaultPageSize: defaultPageSize,
	}
}
