package monitor

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/h-karlsson/uptime/internal/validation"
)

type Store interface {
	List(ctx context.Context, userID int64) ([]Monitor, error)
	Get(ctx context.Context, userID, id int64) (Monitor, error)
	Count(ctx context.Context, userID int64) (int, error)
	Create(ctx context.Context, userID int64, input Input) (Monitor, error)
	Update(ctx context.Context, userID, id int64, input Input) (Monitor, error)
	Delete(ctx context.Context, userID, id int64) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) List(ctx context.Context, userID int64) ([]Monitor, error) {
	return service.store.List(ctx, userID)
}

func (service *Service) Get(ctx context.Context, userID, id int64) (Monitor, error) {
	return service.store.Get(ctx, userID, id)
}

func (service *Service) Create(ctx context.Context, userID int64, input Input) (Monitor, error) {
	input, err := normalize(input)
	if err != nil {
		return Monitor{}, err
	}

	// Count-then-insert can race, but the limit only needs to stop runaway use, not be exact.
	count, err := service.store.Count(ctx, userID)
	if err != nil {
		return Monitor{}, err
	}
	if count >= MaxPerUser {
		return Monitor{}, ErrLimitReached
	}
	return service.store.Create(ctx, userID, input)
}

func (service *Service) Update(ctx context.Context, userID, id int64, input Input) (Monitor, error) {
	input, err := normalize(input)
	if err != nil {
		return Monitor{}, err
	}
	return service.store.Update(ctx, userID, id, input)
}

func (service *Service) Delete(ctx context.Context, userID, id int64) error {
	return service.store.Delete(ctx, userID, id)
}

func normalize(input Input) (Input, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Method = strings.ToUpper(strings.TrimSpace(input.Method))
	canonical, validURL := canonicalURL(input.URL)
	input.URL = canonical

	switch {
	case input.Name == "":
		return Input{}, invalid("name is required")
	case len(input.Name) > MaxNameLength:
		return Input{}, invalid("name must be at most %d characters", MaxNameLength)
	case !validURL:
		return Input{}, invalid("url must be an http or https address")
	case len(input.URL) > MaxURLLength:
		return Input{}, invalid("url must be at most %d characters", MaxURLLength)
	case input.Method != "GET" && input.Method != "HEAD":
		return Input{}, invalid("method must be GET or HEAD")
	case input.IntervalSeconds < MinIntervalSeconds || input.IntervalSeconds > MaxIntervalSeconds:
		return Input{}, invalid("interval must be between %d and %d seconds", MinIntervalSeconds, MaxIntervalSeconds)
	case input.TimeoutMs < MinTimeoutMs || input.TimeoutMs > MaxTimeoutMs:
		return Input{}, invalid("timeout must be between %d and %d ms", MinTimeoutMs, MaxTimeoutMs)
	case input.ExpectedStatus < 100 || input.ExpectedStatus > 599:
		return Input{}, invalid("expected status must be a valid HTTP status code")
	}
	return input, nil
}

// canonicalURL lowercases the host, adds the root path, and drops the fragment, so
// equivalent addresses like "https://Example.com" and "https://example.com/" are duplicates.
func canonicalURL(raw string) (string, bool) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", false
	}
	parsed.Host = strings.ToLower(parsed.Host)
	parsed.Fragment = ""
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return parsed.String(), true
}

func invalid(format string, args ...any) error {
	return &validation.Error{Message: fmt.Sprintf(format, args...)}
}
