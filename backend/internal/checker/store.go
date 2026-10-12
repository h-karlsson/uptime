package checker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/h-karlsson/uptime/internal/store/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A monitor is only marked down after this many failed checks in a row, so a single slow
// response doesn't raise a false alarm.
const failureThreshold = 3

type Store struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, queries: db.New(pool)}
}

func (store *Store) ClaimDue(ctx context.Context, batchSize int) ([]target, error) {
	rows, err := store.queries.ClaimDueMonitors(ctx, int32(batchSize))
	if err != nil {
		return nil, fmt.Errorf("claim due monitors: %w", err)
	}
	targets := make([]target, len(rows))
	for i, row := range rows {
		targets[i] = target{
			monitorID:      row.ID,
			url:            row.Url,
			method:         row.Method,
			timeout:        time.Duration(row.TimeoutMs) * time.Millisecond,
			expectedStatus: int(row.ExpectedStatus),
		}
	}
	return targets, nil
}

// Record stores a check result and updates the monitor's status, opening an incident when it
// goes down and resolving it when it recovers, all in one transaction.
func (store *Store) Record(ctx context.Context, monitorID int64, result result) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	queries := store.queries.WithTx(tx)

	health, err := queries.UpdateMonitorHealth(ctx, db.UpdateMonitorHealthParams{
		Ok:               result.ok,
		FailureThreshold: failureThreshold,
		ID:               monitorID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil // deleted while it was being checked
	}
	if err != nil {
		return fmt.Errorf("update monitor health: %w", err)
	}

	err = queries.InsertCheck(ctx, db.InsertCheckParams{
		MonitorID:  monitorID,
		CheckedAt:  result.checkedAt,
		Ok:         result.ok,
		StatusCode: result.statusCode,
		LatencyMs:  result.latencyMs,
		Error:      result.err,
	})
	if err != nil {
		return fmt.Errorf("insert check: %w", err)
	}

	switch {
	case health.CurrentStatus == "down" && health.PreviousStatus != "down":
		slog.Info("monitor down", "monitor", monitorID, "cause", *result.err)
		err = queries.OpenIncident(ctx, db.OpenIncidentParams{MonitorID: monitorID, StartedAt: result.checkedAt, Cause: result.err})
	case health.CurrentStatus == "up" && health.PreviousStatus == "down":
		slog.Info("monitor recovered", "monitor", monitorID)
		err = queries.ResolveIncident(ctx, db.ResolveIncidentParams{MonitorID: monitorID, ResolvedAt: &result.checkedAt})
	}
	if err != nil {
		return fmt.Errorf("update incident: %w", err)
	}
	return tx.Commit(ctx)
}
