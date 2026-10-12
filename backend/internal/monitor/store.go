package monitor

import (
	"context"
	"errors"
	"fmt"

	"github.com/h-karlsson/uptime/internal/pgerr"
	"github.com/h-karlsson/uptime/internal/store/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueURLIndex = "monitors_user_id_url_idx"

// PostgresStore scopes every query to the given user, so a monitor owned by someone else
// behaves exactly like one that doesn't exist.
type PostgresStore struct {
	queries *db.Queries
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{queries: db.New(pool)}
}

func (store *PostgresStore) List(ctx context.Context, userID int64) ([]Monitor, error) {
	rows, err := store.queries.ListMonitors(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list monitors: %w", err)
	}
	monitors := make([]Monitor, len(rows))
	for i, row := range rows {
		monitors[i] = fromRow(row)
	}
	return monitors, nil
}

func (store *PostgresStore) Get(ctx context.Context, userID, id int64) (Monitor, error) {
	row, err := store.queries.GetMonitor(ctx, db.GetMonitorParams{ID: id, UserID: userID})
	if errors.Is(err, pgx.ErrNoRows) {
		return Monitor{}, ErrNotFound
	}
	if err != nil {
		return Monitor{}, fmt.Errorf("get monitor: %w", err)
	}
	return fromRow(row), nil
}

func (store *PostgresStore) Count(ctx context.Context, userID int64) (int, error) {
	count, err := store.queries.CountMonitors(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("count monitors: %w", err)
	}
	return int(count), nil
}

func (store *PostgresStore) Create(ctx context.Context, userID int64, input Input) (Monitor, error) {
	row, err := store.queries.CreateMonitor(ctx, db.CreateMonitorParams{
		UserID:          userID,
		Name:            input.Name,
		Url:             input.URL,
		Method:          input.Method,
		IntervalSeconds: int32(input.IntervalSeconds),
		TimeoutMs:       int32(input.TimeoutMs),
		ExpectedStatus:  int32(input.ExpectedStatus),
		Enabled:         input.Enabled,
	})
	if pgerr.IsUniqueViolation(err, uniqueURLIndex) {
		return Monitor{}, ErrDuplicateURL
	}
	if err != nil {
		return Monitor{}, fmt.Errorf("insert monitor: %w", err)
	}
	return fromRow(row), nil
}

func (store *PostgresStore) Update(ctx context.Context, userID, id int64, input Input) (Monitor, error) {
	row, err := store.queries.UpdateMonitor(ctx, db.UpdateMonitorParams{
		ID:              id,
		UserID:          userID,
		Name:            input.Name,
		Url:             input.URL,
		Method:          input.Method,
		IntervalSeconds: int32(input.IntervalSeconds),
		TimeoutMs:       int32(input.TimeoutMs),
		ExpectedStatus:  int32(input.ExpectedStatus),
		Enabled:         input.Enabled,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return Monitor{}, ErrNotFound
	}
	if pgerr.IsUniqueViolation(err, uniqueURLIndex) {
		return Monitor{}, ErrDuplicateURL
	}
	if err != nil {
		return Monitor{}, fmt.Errorf("update monitor: %w", err)
	}
	return fromRow(row), nil
}

func (store *PostgresStore) Delete(ctx context.Context, userID, id int64) error {
	deleted, err := store.queries.DeleteMonitor(ctx, db.DeleteMonitorParams{ID: id, UserID: userID})
	if err != nil {
		return fmt.Errorf("delete monitor: %w", err)
	}
	if deleted == 0 {
		return ErrNotFound
	}
	return nil
}

func fromRow(row db.Monitor) Monitor {
	return Monitor{
		ID:              row.ID,
		Name:            row.Name,
		URL:             row.Url,
		Method:          row.Method,
		IntervalSeconds: int(row.IntervalSeconds),
		TimeoutMs:       int(row.TimeoutMs),
		ExpectedStatus:  int(row.ExpectedStatus),
		Enabled:         row.Enabled,
		Status:          row.CurrentStatus,
		CreatedAt:       row.CreatedAt,
	}
}
