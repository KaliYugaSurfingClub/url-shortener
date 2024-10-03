package linkManager

import (
	"context"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/utils"
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

func (m *LinkManager) GetUsersLinks(ctx context.Context, userId int64, params model.GetLinksParams) (links []*model.Link, totalCount int64, err error) {
	defer utils.WithinOp("core.linkManager.GetUsersLinks", &err)

	if totalCount, err = m.links.GetCountByUserId(ctx, userId, params.Filter); err != nil {
		return nil, 0, err
	}

	if links, err = m.links.GetByUserId(ctx, userId, params); err != nil {
		return nil, 0, err
	}

	return links, totalCount, nil
}

func (m *LinkManager) GetLinkWithClicks(ctx context.Context, linkId int64, params model.GetClicksParams) (link *model.Link, clicks []*model.Click, totalCount int64, err error) {
	defer utils.WithinOp("core.linkManager.GetLinkWithClicks", &err)

	if link, err = m.links.GetById(ctx, linkId); err != nil {
		return nil, nil, 0, err
	}

	if totalCount, err = m.clicks.GetCountByLinkId(ctx, linkId, params); err != nil {
		return nil, nil, 0, err
	}

	if clicks, err = m.clicks.GetByLinkId(ctx, linkId, params); err != nil {
		return nil, nil, 0, err
	}

	return link, clicks, totalCount, nil
}
