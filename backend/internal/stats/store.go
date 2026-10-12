package stats

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/h-karlsson/uptime/internal/store/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	queries *db.Queries
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{queries: db.New(pool)}
}

func (store *PostgresStore) BelongsToUser(ctx context.Context, userID, monitorID int64) (bool, error) {
	owned, err := store.queries.MonitorBelongsToUser(ctx, db.MonitorBelongsToUserParams{ID: monitorID, UserID: userID})
	if err != nil {
		return false, fmt.Errorf("check monitor owner: %w", err)
	}
	return owned, nil
}

func (store *PostgresStore) Summary(ctx context.Context, monitorID int64, since time.Time) (Summary, error) {
	row, err := store.queries.GetCheckSummary(ctx, db.GetCheckSummaryParams{MonitorID: monitorID, Since: since})
	if err != nil {
		return Summary{}, fmt.Errorf("summarize checks: %w", err)
	}
	summary := Summary{
		Checks:        int(row.Checks),
		UptimePercent: uptimePercent(row.OkChecks, row.Checks),
		AvgLatencyMs:  latency(row.AvgLatencyMs, row.OkChecks),
		P95LatencyMs:  latency(row.P95LatencyMs, row.OkChecks),
	}

	lastCheckedAt, err := store.queries.GetLastCheckedAt(ctx, monitorID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return Summary{}, fmt.Errorf("get last check: %w", err)
	}
	if err == nil {
		summary.LastCheckedAt = &lastCheckedAt
	}
	return summary, nil
}

func (store *PostgresStore) Latency(ctx context.Context, monitorID int64, since time.Time, bucket time.Duration) ([]LatencyPoint, error) {
	rows, err := store.queries.ListLatencyBuckets(ctx, db.ListLatencyBucketsParams{
		BucketSeconds: int32(bucket.Seconds()),
		MonitorID:     monitorID,
		Since:         since,
	})
	if err != nil {
		return nil, fmt.Errorf("list latency buckets: %w", err)
	}
	points := make([]LatencyPoint, len(rows))
	for i, row := range rows {
		points[i] = LatencyPoint{
			Time:         row.BucketStart,
			Checks:       int(row.Checks),
			FailedChecks: int(row.Checks - row.OkChecks),
			AvgLatencyMs: latency(row.AvgLatencyMs, row.OkChecks),
		}
	}
	return points, nil
}

func (store *PostgresStore) Days(ctx context.Context, monitorID int64, since time.Time) ([]Day, error) {
	rows, err := store.queries.ListDailyUptime(ctx, db.ListDailyUptimeParams{MonitorID: monitorID, Since: since})
	if err != nil {
		return nil, fmt.Errorf("list daily uptime: %w", err)
	}
	days := make([]Day, len(rows))
	for i, row := range rows {
		days[i] = Day{
			Date:          row.Day.Format(time.DateOnly),
			Checks:        int(row.Checks),
			UptimePercent: *uptimePercent(row.OkChecks, row.Checks),
		}
	}
	return days, nil
}

func (store *PostgresStore) Incidents(ctx context.Context, monitorID int64) ([]Incident, error) {
	rows, err := store.queries.ListIncidents(ctx, monitorID)
	if err != nil {
		return nil, fmt.Errorf("list incidents: %w", err)
	}
	incidents := make([]Incident, len(rows))
	for i, row := range rows {
		incidents[i] = Incident{ID: row.ID, StartedAt: row.StartedAt, ResolvedAt: row.ResolvedAt, Cause: row.Cause}
	}
	return incidents, nil
}

func uptimePercent(ok, total int32) *float64 {
	if total == 0 {
		return nil
	}
	percent := float64(ok) / float64(total) * 100
	return &percent
}

// latency returns nil when there were no successful checks; the queries report 0 in that case
// because sqlc can't type aggregates as nullable.
func latency(ms float64, okChecks int32) *int {
	if okChecks == 0 {
		return nil
	}
	rounded := int(math.Round(ms))
	return &rounded
}
