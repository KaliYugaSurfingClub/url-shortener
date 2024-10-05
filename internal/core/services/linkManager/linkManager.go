package linkManager

import (
	"context"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/utils"
)

type LinkManager struct {
	links port.LinkStorage
}

func New(linkStorage port.LinkStorage) *LinkManager {
	return &LinkManager{
		links: linkStorage,
	}
}

// todo duplicate 1
// todo userId in params params is pointer
func (m *LinkManager) GetUserLinks(ctx context.Context, userId int64, params model.GetLinksParams) (links []*model.Link, totalCount int64, err error) {
	defer utils.WithinOp("core.linkManager.GetUserLinks", &err)

	if totalCount, err = m.links.GetCountByUserId(ctx, userId, params.Filter); err != nil {
		return nil, 0, err
	}

	if links, err = m.links.GetByUserId(ctx, userId, params); err != nil {
		return nil, 0, err
	}

	return links, totalCount, nil
}
