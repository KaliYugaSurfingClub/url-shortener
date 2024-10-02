package linkManager

import (
	"context"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type LinkManager struct {
	links  port.LinkStorage
	clicks port.ClickStorage
}

func New(linkStorage port.LinkStorage, clickStorage port.ClickStorage) *LinkManager {
	return &LinkManager{
		links:  linkStorage,
		clicks: clickStorage,
	}
}

func (m *LinkManager) GetUsersLinks(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, int64, error) {
	totalCount, err := m.links.GetCountByUserId(ctx, userId, params.Filter)
	if err != nil {
		return nil, 0, err
	}

	links, err := m.links.GetByUserId(ctx, userId, params)
	if err != nil {
		return nil, 0, err
	}

	return links, totalCount, nil
}

func (m *LinkManager) GetLinkById(ctx context.Context, linkId int64) {

}
