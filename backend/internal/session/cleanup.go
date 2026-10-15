package session

import (
	"context"
	"log/slog"
	"time"
)

func (service *Service) RunCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		deleted, err := service.store.DeleteExpired(ctx)
		switch {
		case err != nil && ctx.Err() == nil:
			slog.Error("clean up sessions", "error", err)
		case deleted > 0:
			slog.Info("cleaned up expired sessions", "deleted", deleted)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
