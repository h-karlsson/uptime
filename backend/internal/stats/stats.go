package stats

import (
	"errors"
	"time"

	"github.com/h-karlsson/uptime/internal/validation"
)

// Range is a reporting window and the bucket size that keeps its latency chart around 100 points.
type Range struct {
	Name     string
	Duration time.Duration
	Bucket   time.Duration
}

var ranges = []Range{
	{Name: "24h", Duration: 24 * time.Hour, Bucket: 15 * time.Minute},
	{Name: "7d", Duration: 7 * 24 * time.Hour, Bucket: 2 * time.Hour},
	{Name: "30d", Duration: 30 * 24 * time.Hour, Bucket: 6 * time.Hour},
}

// ParseRange returns the named range, or 24h when name is empty.
func ParseRange(name string) (Range, error) {
	if name == "" {
		return ranges[0], nil
	}
	for _, candidate := range ranges {
		if candidate.Name == name {
			return candidate, nil
		}
	}
	return Range{}, &validation.Error{Message: "range must be 24h, 7d or 30d"}
}

// DailyHistory is how many days of uptime the status bars show, including today.
const DailyHistory = 30

type Summary struct {
	Range         string         `json:"range"`
	Checks        int            `json:"checks"`
	UptimePercent *float64       `json:"uptimePercent"`
	AvgLatencyMs  *int           `json:"avgLatencyMs"`
	P95LatencyMs  *int           `json:"p95LatencyMs"`
	LastCheckedAt *time.Time     `json:"lastCheckedAt"`
	Latency       []LatencyPoint `json:"latency"`
}

type LatencyPoint struct {
	Time         time.Time `json:"time"`
	Checks       int       `json:"checks"`
	FailedChecks int       `json:"failedChecks"`
	AvgLatencyMs *int      `json:"avgLatencyMs"`
}

type Day struct {
	Date          string  `json:"date"`
	Checks        int     `json:"checks"`
	UptimePercent float64 `json:"uptimePercent"`
}

type Incident struct {
	ID         int64      `json:"id"`
	StartedAt  time.Time  `json:"startedAt"`
	ResolvedAt *time.Time `json:"resolvedAt"`
	Cause      *string    `json:"cause"`
}

var ErrNotFound = errors.New("monitor not found")
