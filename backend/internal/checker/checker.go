package checker

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const partitionCheckInterval = time.Hour

type Config struct {
	Workers      int
	PollInterval time.Duration
}

// Checker claims due monitors from Postgres and checks them with a pool of workers.
type Checker struct {
	store  *Store
	client *http.Client
	config Config
}

func New(pool *pgxpool.Pool, config Config) *Checker {
	return &Checker{store: NewStore(pool), client: newClient(), config: config}
}

// Run checks monitors until ctx is cancelled. It then stops claiming new work and waits for
// the checks already in progress, so their results are still recorded.
func (checker *Checker) Run(ctx context.Context) {
	targets := make(chan target)
	var workers sync.WaitGroup
	for range checker.config.Workers {
		workers.Go(func() {
			for target := range targets {
				checker.check(context.WithoutCancel(ctx), target)
			}
		})
	}
	defer workers.Wait()
	defer close(targets)

	slog.Info("checker started", "workers", checker.config.Workers)
	ticker := time.NewTicker(checker.config.PollInterval)
	defer ticker.Stop()

	var partitionsCheckedAt time.Time
	for {
		if time.Since(partitionsCheckedAt) >= partitionCheckInterval {
			if err := checker.store.EnsurePartitions(ctx, time.Now()); err != nil {
				slog.Error("ensure checks partitions", "error", err)
			} else {
				partitionsCheckedAt = time.Now()
			}
		}
		checker.dispatchDue(ctx, targets)

		select {
		case <-ctx.Done():
			slog.Info("checker stopping")
			return
		case <-ticker.C:
		}
	}
}

// dispatchDue hands due monitors to the workers until none are left. Claiming at most one
// batch per free worker keeps monitors from being claimed long before a worker can check them.
func (checker *Checker) dispatchDue(ctx context.Context, targets chan<- target) {
	for ctx.Err() == nil {
		due, err := checker.store.ClaimDue(ctx, checker.config.Workers)
		if err != nil {
			slog.Error("claim due monitors", "error", err)
			return
		}
		for _, target := range due {
			select {
			case targets <- target:
			case <-ctx.Done():
				return
			}
		}
		if len(due) < checker.config.Workers {
			return
		}
	}
}

func (checker *Checker) check(ctx context.Context, target target) {
	result := probe(ctx, checker.client, target)
	if err := checker.store.Record(ctx, target.monitorID, result); err != nil {
		slog.Error("record check", "monitor", target.monitorID, "error", err)
	}
}
