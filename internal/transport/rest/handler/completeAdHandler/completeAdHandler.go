package completeAdHandler

import "context"

type AdCompleter interface {
	CompleteAd(ctx context.Context, clickId int64)
}

type Handler struct {
}
