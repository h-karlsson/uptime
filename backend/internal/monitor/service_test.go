package monitor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/h-karlsson/uptime/internal/validation"
)

func TestCreateValidatesInput(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Input)
	}{
		{"blank name", func(input *Input) { input.Name = "   " }},
		{"long name", func(input *Input) { input.Name = strings.Repeat("a", MaxNameLength+1) }},
		{"not a url", func(input *Input) { input.URL = "example.com" }},
		{"unsupported scheme", func(input *Input) { input.URL = "ftp://example.com" }},
		{"missing host", func(input *Input) { input.URL = "https://" }},
		{"long url", func(input *Input) { input.URL = "https://example.com/" + strings.Repeat("a", MaxURLLength) }},
		{"unsupported method", func(input *Input) { input.Method = "POST" }},
		{"interval too short", func(input *Input) { input.IntervalSeconds = MinIntervalSeconds - 1 }},
		{"interval too long", func(input *Input) { input.IntervalSeconds = MaxIntervalSeconds + 1 }},
		{"timeout too short", func(input *Input) { input.TimeoutMs = MinTimeoutMs - 1 }},
		{"timeout too long", func(input *Input) { input.TimeoutMs = MaxTimeoutMs + 1 }},
		{"invalid status", func(input *Input) { input.ExpectedStatus = 99 }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := validInput()
			test.change(&input)
			_, err := NewService(&fakeStore{}).Create(context.Background(), 1, input)
			if _, ok := errors.AsType[*validation.Error](err); !ok {
				t.Fatalf("err = %v, want a validation error", err)
			}
		})
	}
}

func TestCreateNormalizesInput(t *testing.T) {
	input := validInput()
	input.Name = "  Example  "
	input.URL = " HTTPS://Example.COM/Health#top "
	input.Method = "head"

	created, err := NewService(&fakeStore{}).Create(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.Name != "Example" || created.URL != "https://example.com/Health" || created.Method != "HEAD" {
		t.Fatalf("got %+v, want trimmed name, canonical url and uppercase method", created)
	}
}

func TestDuplicateURLs(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)
	ctx := context.Background()

	first, err := service.Create(ctx, 1, validInput())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	equivalent := validInput()
	equivalent.URL = "HTTPS://EXAMPLE.com/"
	if _, err := service.Create(ctx, 1, equivalent); !errors.Is(err, ErrDuplicateURL) {
		t.Fatalf("equivalent url: err = %v, want ErrDuplicateURL", err)
	}

	other := validInput()
	other.URL = "https://example.org"
	second, err := service.Create(ctx, 1, other)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := service.Update(ctx, 1, second.ID, validInput()); !errors.Is(err, ErrDuplicateURL) {
		t.Fatalf("update to a taken url: err = %v, want ErrDuplicateURL", err)
	}

	renamed := validInput()
	renamed.Name = "Renamed"
	if _, err := service.Update(ctx, 1, first.ID, renamed); err != nil {
		t.Fatalf("updating a monitor without changing its url: %v", err)
	}
	if _, err := service.Create(ctx, 2, validInput()); err != nil {
		t.Fatalf("another user can monitor the same url: %v", err)
	}
}

func TestCreateEnforcesLimitPerUser(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)
	ctx := context.Background()

	inputFor := func(n int) Input {
		input := validInput()
		input.URL = fmt.Sprintf("https://example.com/%d", n)
		return input
	}
	for n := range MaxPerUser {
		if _, err := service.Create(ctx, 1, inputFor(n)); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}
	if _, err := service.Create(ctx, 1, inputFor(MaxPerUser)); !errors.Is(err, ErrLimitReached) {
		t.Fatalf("err = %v, want ErrLimitReached", err)
	}
	if _, err := service.Create(ctx, 2, validInput()); err != nil {
		t.Fatalf("another user's limit is separate: %v", err)
	}
}
