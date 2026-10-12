package checker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

const userAgent = "UptimeMonitor/0.1 (+https://github.com/h-karlsson/uptime)"

type target struct {
	monitorID      int64
	url            string
	method         string
	timeout        time.Duration
	expectedStatus int
}

type result struct {
	checkedAt  time.Time
	ok         bool
	statusCode *int32
	latencyMs  *int32
	err        *string
}

// probe sends one request and reports whether the target answered with the expected status.
func probe(ctx context.Context, client *http.Client, target target) result {
	checkedAt := time.Now()
	ctx, cancel := context.WithTimeout(ctx, target.timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, target.method, target.url, nil)
	if err != nil {
		return failed(checkedAt, err.Error())
	}
	request.Header.Set("User-Agent", userAgent)

	response, err := client.Do(request)
	if err != nil {
		return failed(checkedAt, describe(err, target.timeout))
	}
	defer response.Body.Close()

	latency := int32(time.Since(checkedAt).Milliseconds())
	status := int32(response.StatusCode)
	result := result{checkedAt: checkedAt, ok: response.StatusCode == target.expectedStatus, statusCode: &status, latencyMs: &latency}
	if !result.ok {
		message := fmt.Sprintf("expected status %d, got %d", target.expectedStatus, response.StatusCode)
		result.err = &message
	}
	return result
}

func failed(checkedAt time.Time, message string) result {
	return result{checkedAt: checkedAt, err: &message}
}

// describe turns a client error into a short message without the request URL that
// net/http prefixes it with.
func describe(err error, timeout time.Duration) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Sprintf("timed out after %s", timeout)
	}
	if errors.Is(err, errBlockedAddress) {
		return "blocked: " + errBlockedAddress.Error()
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return "could not resolve host " + dnsErr.Name
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return urlErr.Err.Error()
	}
	return err.Error()
}
