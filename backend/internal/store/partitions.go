package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureCheckPartitions creates the monthly checks partitions for the month of now and the next.
func EnsureCheckPartitions(ctx context.Context, pool *pgxpool.Pool, now time.Time) error {
	month := time.Date(now.UTC().Year(), now.UTC().Month(), 1, 0, 0, 0, 0, time.UTC)
	for _, start := range []time.Time{month, month.AddDate(0, 1, 0)} {
		end := start.AddDate(0, 1, 0)
		statement := fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS checks_%s PARTITION OF checks FOR VALUES FROM ('%s') TO ('%s')",
			start.Format("2006_01"), start.Format(time.RFC3339), end.Format(time.RFC3339),
		)
		if _, err := pool.Exec(ctx, statement); err != nil {
			return fmt.Errorf("create checks partition for %s: %w", start.Format("2006-01"), err)
		}
	}
	return nil
}
