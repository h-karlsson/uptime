package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/testdb"
)

func TestPostgresStore(t *testing.T) {
	pool := testdb.Open(t)
	store := NewPostgresStore(pool)
	ctx := context.Background()
	alice := testdb.CreateUser(t, pool, "alice@uptime.io")
	bob := testdb.CreateUser(t, pool, "bob@uptime.io")

	create := func(tokenHash string, userID int64, expiresAt time.Time) {
		t.Helper()
		if err := store.Create(ctx, []byte(tokenHash), userID, expiresAt); err != nil {
			t.Fatal(err)
		}
	}
	valid := func(tokenHash string) bool {
		t.Helper()
		_, err := store.Get(ctx, []byte(tokenHash))
		if err != nil && !errors.Is(err, ErrInvalidSession) {
			t.Fatal(err)
		}
		return err == nil
	}

	inAnHour := time.Now().Add(time.Hour)
	create("alice-laptop", alice, inAnHour)
	create("alice-phone", alice, inAnHour)
	create("alice-expired", alice, time.Now().Add(-time.Hour))
	create("bob", bob, inAnHour)

	if valid("alice-expired") {
		t.Fatal("expired session should be rejected")
	}
	current, err := store.Get(ctx, []byte("alice-laptop"))
	if err != nil || current.User.Email != "alice@uptime.io" || !current.ExpiresAt.Round(time.Millisecond).Equal(inAnHour.Round(time.Millisecond)) {
		t.Fatalf("Get = %+v, %v; want alice expiring at %v", current, err, inAnHour)
	}

	extended := time.Now().Add(Duration)
	if err := store.Extend(ctx, []byte("alice-laptop"), extended); err != nil {
		t.Fatal(err)
	}
	if current, _ := store.Get(ctx, []byte("alice-laptop")); current.ExpiresAt.Before(extended.Add(-time.Second)) {
		t.Fatalf("after Extend expiresAt = %v, want %v", current.ExpiresAt, extended)
	}

	deleted, err := store.DeleteExpired(ctx)
	if err != nil || deleted != 1 {
		t.Fatalf("DeleteExpired = %d, %v; want 1", deleted, err)
	}

	if err := store.DeleteForUser(ctx, alice); err != nil {
		t.Fatal(err)
	}
	if valid("alice-laptop") || valid("alice-phone") {
		t.Fatal("DeleteForUser left alice's sessions")
	}
	if !valid("bob") {
		t.Fatal("DeleteForUser removed another user's session")
	}
}
