package monitor

import (
	"errors"
	"time"
)

type Monitor struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	URL             string    `json:"url"`
	Method          string    `json:"method"`
	IntervalSeconds int       `json:"intervalSeconds"`
	TimeoutMs       int       `json:"timeoutMs"`
	ExpectedStatus  int       `json:"expectedStatus"`
	Enabled         bool      `json:"enabled"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Input struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Method          string `json:"method"`
	IntervalSeconds int    `json:"intervalSeconds"`
	TimeoutMs       int    `json:"timeoutMs"`
	ExpectedStatus  int    `json:"expectedStatus"`
	Enabled         bool   `json:"enabled"`
}

const (
	MaxPerUser         = 10
	MaxNameLength      = 100
	MaxURLLength       = 2048
	MinIntervalSeconds = 30
	MaxIntervalSeconds = 24 * 60 * 60
	MinTimeoutMs       = 1000
	MaxTimeoutMs       = 30000
)

var (
	ErrNotFound     = errors.New("monitor not found")
	ErrLimitReached = errors.New("monitor limit reached")
	ErrDuplicateURL = errors.New("you already have a monitor for this url")
)
