package clickManager

import (
	"context"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/utils"
)

type LinkManager struct {
	clicks port.ClickStorage
}

func New(clickStorage port.ClickStorage) *LinkManager {
	return &LinkManager{
		clicks: clickStorage,
	}
}

// todo duplicate 1
func (m *LinkManager) GetLinkClicks(ctx context.Context, linkId int64, params model.GetClicksParams) (clicks []*model.Click, totalCount int64, err error) {
	defer utils.WithinOp("core.linkManager.GetUserLinks", &err)

	if totalCount, err = m.clicks.GetCountByLinkId(ctx, linkId, params); err != nil {
		return nil, 0, err
	}

	if clicks, err = m.clicks.GetByLinkId(ctx, linkId, params); err != nil {
		return nil, 0, err
	}

	return clicks, totalCount, nil
}
