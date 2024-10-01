package linkManager

import (
	"context"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type LinkManager struct {
	provider port.LinkStorage
}

func (m *LinkManager) GetUsersLinks(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, int64, error) {
	totalCount, err := m.provider.GetCountByUserId(ctx, userId, params.Filter)
	if err != nil {
		return nil, 0, err
	}

	links, err := m.provider.GetByUserId(ctx, userId, params)
	if err != nil {
		return nil, 0, err
	}

	return links, totalCount, nil
}
