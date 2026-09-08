package legacy

import (
	"context"
	"testing"
)

func TestStorePutGet(t *testing.T) {
	ctx := context.Background()

	s := NewStore()
	t.Cleanup(func() {
		if err := s.Flush(ctx); err != nil {
			t.Errorf("flush: %v", err)
		}
	})

	if err := s.Put(ctx, "a", "1"); err != nil {
		t.Fatal(err)
	}
	got, err := s.Get(ctx, "a")
	if err != nil {
		t.Fatal(err)
	}
	if got != "1" {
		t.Fatalf("got %q, want %q", got, "1")
	}
}
