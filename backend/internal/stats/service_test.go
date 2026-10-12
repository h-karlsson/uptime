package stats

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/validation"
)

func TestSummaryUsesRangeWindowAndBucket(t *testing.T) {
	tests := []struct {
		rangeName string
		want      string
		duration  time.Duration
		bucket    time.Duration
	}{
		{"", "24h", 24 * time.Hour, 15 * time.Minute},
		{"24h", "24h", 24 * time.Hour, 15 * time.Minute},
		{"7d", "7d", 7 * 24 * time.Hour, 2 * time.Hour},
		{"30d", "30d", 30 * 24 * time.Hour, 6 * time.Hour},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			store := &fakeStore{}
			summary, err := NewService(store).Summary(context.Background(), 1, 1, test.rangeName)
			if err != nil {
				t.Fatal(err)
			}
			if summary.Range != test.want || store.bucket != test.bucket {
				t.Fatalf("range %q with bucket %s, want %q with %s", summary.Range, store.bucket, test.want, test.bucket)
			}
			if window := time.Since(store.since); window < test.duration || window > test.duration+time.Minute {
				t.Fatalf("window = %s, want %s", window, test.duration)
			}
		})
	}
}

func TestSummaryRejectsUnknownRange(t *testing.T) {
	_, err := NewService(&fakeStore{}).Summary(context.Background(), 1, 1, "1y")
	if _, ok := errors.AsType[*validation.Error](err); !ok {
		t.Fatalf("err = %v, want a validation error", err)
	}
}

func TestDaysStartAtMidnightUTC(t *testing.T) {
	store := &fakeStore{}
	if _, err := NewService(store).Days(context.Background(), 1, 1); err != nil {
		t.Fatal(err)
	}
	want := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, -(DailyHistory - 1))
	if !store.since.Equal(want) {
		t.Fatalf("since = %s, want %s", store.since, want)
	}
}

func TestOtherUsersMonitorIsNotFound(t *testing.T) {
	service := NewService(&fakeStore{})
	ctx := context.Background()

	if _, err := service.Summary(ctx, 2, 1, "24h"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("summary: err = %v, want ErrNotFound", err)
	}
	if _, err := service.Days(ctx, 2, 1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("days: err = %v, want ErrNotFound", err)
	}
	if _, err := service.Incidents(ctx, 2, 1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("incidents: err = %v, want ErrNotFound", err)
	}
}
