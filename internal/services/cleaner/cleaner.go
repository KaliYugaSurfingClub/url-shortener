package cleaner

import (
	"log/slog"
	"time"
)

//todo delete alias-related information from other tables using transactions

type Deleter interface {
	DeleteOverdueAliases(deadline time.Time) (int64, error)
}

func Start(log *slog.Logger, deleter Deleter, lifetime time.Duration, period time.Duration) {
	log.Info("cleaner started")

	for range time.Tick(period) {
		deadline := time.Now().Add(-lifetime)
		count, err := deleter.DeleteOverdueAliases(deadline)

		switch {
		case err != nil:
			log.Error("deleting overdue aliases was failed")
		case count > 0:
			log.Info("deleted overdue aliases", slog.Int64("count", count))
		default:
			log.Debug("no overdue aliases")
		}
	}
}
