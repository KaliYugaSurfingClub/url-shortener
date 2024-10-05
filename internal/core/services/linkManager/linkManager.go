package linkManager

import (
	"context"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/utils"
)

type LinkManager struct {
	links  port.LinkStorage
	clicks port.ClickStorage
}

func New(linkStorage port.LinkStorage, clicksStorage port.ClickStorage) *LinkManager {
	return &LinkManager{
		links:  linkStorage,
		clicks: clicksStorage,
	}
}

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

func (m *LinkManager) GetLinkClicks(ctx context.Context, linkId int64, userId int64, params model.GetClicksParams) (clicks []*model.Click, totalCount int64, err error) {
	defer utils.WithinOp("core.linkManager.GetLinkClicks", &err)

	ok, err := m.links.DoesLinkBelongUser(ctx, linkId, userId)
	if err != nil {
		return nil, 0, err
	}
	if !ok {
		return nil, 0, core.ErrLinkDoesNotBelongsUser
	}

	if totalCount, err = m.clicks.GetCountByLinkId(ctx, linkId, params); err != nil {
		return nil, 0, err
	}

	if clicks, err = m.clicks.GetByLinkId(ctx, linkId, params); err != nil {
		return nil, 0, err
	}

	return clicks, totalCount, nil
}
