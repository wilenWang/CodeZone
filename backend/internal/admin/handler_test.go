package admin

import "testing"

func TestLimitRecentMessages(t *testing.T) {
	if got := normalizeLimit(0); got != 50 {
		t.Fatalf("got %d want 50", got)
	}
	if got := normalizeLimit(500); got != 100 {
		t.Fatalf("got %d want 100", got)
	}
	if got := normalizeLimit(20); got != 20 {
		t.Fatalf("got %d want 20", got)
	}
}
