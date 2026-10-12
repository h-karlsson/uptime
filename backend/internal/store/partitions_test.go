package store_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/store"
	"github.com/h-karlsson/uptime/internal/testdb"
)

func TestEnsureCheckPartitionsIsIdempotent(t *testing.T) {
	pool := testdb.Open(t)
	ctx := context.Background()
	december := time.Date(2030, time.December, 15, 12, 0, 0, 0, time.UTC)

	for range 2 {
		if err := store.EnsureCheckPartitions(ctx, pool, december); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := pool.Query(ctx, "SELECT relname FROM pg_class WHERE relname IN ('checks_2030_12', 'checks_2031_01') ORDER BY relname")
	if err != nil {
		t.Fatal(err)
	}
	var partitions []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		partitions = append(partitions, name)
	}
	if !slices.Equal(partitions, []string{"checks_2030_12", "checks_2031_01"}) {
		t.Fatalf("partitions = %v, want December and the following January", partitions)
	}
}
