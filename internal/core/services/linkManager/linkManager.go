package linkManager

import (
	"context"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/utils"
)

type LinkManager struct {
	storage port.Repository
}

func New(storage port.Repository) *LinkManager {
	return &LinkManager{
		storage: storage,
	}
}

func (m *LinkManager) GetUserLinks(ctx context.Context, params model.GetLinksParams) (links []*model.Link, totalCount int64, err error) {
	defer utils.WithinOp("core.linkManager.GetUserLinks", &err)

	if totalCount, err = m.storage.GetLinksCountByParams(ctx, params); err != nil {
		return nil, 0, err
	}

	if links, err = m.storage.GetLinksByParams(ctx, params); err != nil {
		return nil, 0, err
	}

	return links, totalCount, nil
}

func (m *LinkManager) GetLinkClicks(ctx context.Context, params model.GetClicksParams) (clicks []*model.Click, totalCount int64, err error) {
	defer utils.WithinOp("core.linkManager.GetLinkClicks", &err)

	ok, err := m.storage.DoesLinkBelongsToUser(ctx, params.LinkId, params.UserId)
	if err != nil {
		return nil, 0, err
	}
	if !ok {
		return nil, 0, core.ErrLinkDoesNotBelongsUser
	}

	if totalCount, err = m.storage.GetClicksCountByParams(ctx, params); err != nil {
		return nil, 0, err
	}

	if clicks, err = m.storage.GetClicksByParams(ctx, params); err != nil {
		return nil, 0, err
	}

	return clicks, totalCount, nil
}
