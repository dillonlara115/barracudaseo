package crawler

import (
	"testing"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/utils"
)

func TestDiscoveredCountsDistinctURLsUpToTheLimit(t *testing.T) {
	m := NewManager(&utils.Config{MaxPages: 3, Workers: 1, Timeout: time.Second, UserAgent: "test"})

	if got := m.Discovered(); got != 0 {
		t.Fatalf("new manager: Discovered() = %d, want 0", got)
	}

	m.noteDiscovered("https://example.com/")
	m.noteDiscovered("https://example.com/a")
	m.noteDiscovered("https://example.com/a") // queued again from another page
	if got := m.Discovered(); got != 2 {
		t.Fatalf("Discovered() = %d, want 2 distinct URLs", got)
	}

	m.noteDiscovered("https://example.com/b")
	m.noteDiscovered("https://example.com/c")
	if got := m.Discovered(); got != 3 {
		t.Fatalf("Discovered() = %d, want it capped at the page limit 3", got)
	}
}
