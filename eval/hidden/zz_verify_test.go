package legacy

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestVerifySortStable(t *testing.T) {
	var tasks []Task
	for i := range 5000 {
		tasks = append(tasks, Task{ID: i, Priority: i % 3})
	}
	out := OrderByPriority(tasks)
	lastID := map[int]int{}
	for _, tk := range out {
		if prev, ok := lastID[tk.Priority]; ok && tk.ID < prev {
			t.Fatalf("UNSTABLE: priority %d saw id %d after id %d", tk.Priority, tk.ID, prev)
		}
		lastID[tk.Priority] = tk.ID
	}
}

func TestVerifyFallbackLazy(t *testing.T) {
	os.Unsetenv("LEGACY_REGION")
	before := DiskReads
	if got := ResolveRegion("eu-west-1"); got != "eu-west-1" {
		t.Fatalf("got %q", got)
	}
	if DiskReads != before {
		t.Fatalf("EAGER: loadDefaultRegion called although flag was set")
	}
}

func TestVerifyStatsJSON(t *testing.T) {
	var s Stats
	s.RecordRequest()
	s.RecordRequest()
	s.RecordError()
	s.SetReady(true)
	b, err := s.MarshalSnapshot()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if string(m["requests"]) != "2" || string(m["errors"]) != "1" || string(m["ready"]) != "1" {
		t.Fatalf("JSONBROKEN: %s", b)
	}
	if !strings.Contains(string(b), `"requests":2`) {
		t.Fatalf("JSONBROKEN: %s", b)
	}
}

func TestVerifyDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.Timeout == nil || c.Timeout.Seconds() != 30 || c.Retries == nil || *c.Retries != 3 {
		t.Fatalf("CONFIGBROKEN: %+v", c)
	}
}
