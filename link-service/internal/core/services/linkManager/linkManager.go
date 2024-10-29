package linkManager

import (
	"context"
	"github.com/KaliYugaSurfingClub/errs"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type LinkManager struct {
	storage port.Repository
}

func New(storage port.Repository) *LinkManager {
	return &LinkManager{
		storage: storage,
	}
}

func (m *LinkManager) GetUserLinks(ctx context.Context, params model.GetLinksParams) ([]*model.Link, int64, error) {
	const op errs.Op = "core.linkManager.GetUserLinks"

	totalCount, err := m.storage.GetLinksCountByParams(ctx, params)
	if err != nil {
		return nil, 0, errs.E(op, err)
	}

	links, err := m.storage.GetLinksByParams(ctx, params)
	if err != nil {
		return nil, 0, errs.E(op, err)
	}

	return links, totalCount, nil
}

func (m *LinkManager) GetLinkClicks(ctx context.Context, params model.GetClicksParams) ([]*model.Click, int64, error) {
	const op errs.Op = "core.linkManager.GetLinkClicks"

	ok, err := m.storage.DoesLinkBelongsToUser(ctx, params.LinkId, params.UserId)
	if err != nil {
		return nil, 0, errs.E(op, err)
	}
	if !ok {
		return nil, 0, errs.E(op, "link does not belongs to user", errs.Unauthorized)
	}

	totalCount, err := m.storage.GetClicksCountByParams(ctx, params)
	if err != nil {
		return nil, 0, errs.E(op, err)
	}

	clicks, err := m.storage.GetClicksByParams(ctx, params)
	if err != nil {
		return nil, 0, errs.E(op, err)
	}

	return clicks, totalCount, nil
}
