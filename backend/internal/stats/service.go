package stats

import (
	"context"
	"time"
)

type Store interface {
	BelongsToUser(ctx context.Context, userID, monitorID int64) (bool, error)
	Summary(ctx context.Context, monitorID int64, since time.Time) (Summary, error)
	Latency(ctx context.Context, monitorID int64, since time.Time, bucket time.Duration) ([]LatencyPoint, error)
	Days(ctx context.Context, monitorID int64, since time.Time) ([]Day, error)
	Incidents(ctx context.Context, monitorID int64) ([]Incident, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) Summary(ctx context.Context, userID, monitorID int64, rangeName string) (Summary, error) {
	reportRange, err := ParseRange(rangeName)
	if err != nil {
		return Summary{}, err
	}
	if err := service.authorize(ctx, userID, monitorID); err != nil {
		return Summary{}, err
	}

	since := time.Now().Add(-reportRange.Duration)
	summary, err := service.store.Summary(ctx, monitorID, since)
	if err != nil {
		return Summary{}, err
	}
	summary.Range = reportRange.Name
	summary.Latency, err = service.store.Latency(ctx, monitorID, since, reportRange.Bucket)
	return summary, err
}

// Days returns uptime per UTC day for the last DailyHistory days. Days without checks are left out.
func (service *Service) Days(ctx context.Context, userID, monitorID int64) ([]Day, error) {
	if err := service.authorize(ctx, userID, monitorID); err != nil {
		return nil, err
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	return service.store.Days(ctx, monitorID, today.AddDate(0, 0, -(DailyHistory-1)))
}

func (service *Service) Incidents(ctx context.Context, userID, monitorID int64) ([]Incident, error) {
	if err := service.authorize(ctx, userID, monitorID); err != nil {
		return nil, err
	}
	return service.store.Incidents(ctx, monitorID)
}

func (service *Service) authorize(ctx context.Context, userID, monitorID int64) error {
	owned, err := service.store.BelongsToUser(ctx, userID, monitorID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrNotFound
	}
	return nil
}
