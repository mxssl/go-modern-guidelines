package legacy

import (
	"encoding/json"
	"sync/atomic"
)

// Stats is a snapshot of server counters. It is serialized to JSON for the
// /debug/stats endpoint and consumed by the dashboard.
type Stats struct {
	Requests int64 `json:"requests"`
	Errors   int64 `json:"errors"`
	Ready    int32 `json:"ready"`
}

// RecordRequest increments the request counter.
func (s *Stats) RecordRequest() {
	atomic.AddInt64(&s.Requests, 1)
}

// RecordError increments the error counter.
func (s *Stats) RecordError() {
	atomic.AddInt64(&s.Errors, 1)
}

// SetReady marks the server ready.
func (s *Stats) SetReady(ready bool) {
	if ready {
		atomic.StoreInt32(&s.Ready, 1)
	} else {
		atomic.StoreInt32(&s.Ready, 0)
	}
}

// IsReady reports readiness.
func (s *Stats) IsReady() bool {
	return atomic.LoadInt32(&s.Ready) == 1
}

// MarshalSnapshot returns the JSON that the dashboard expects.
func (s *Stats) MarshalSnapshot() ([]byte, error) {
	snap := Stats{
		Requests: atomic.LoadInt64(&s.Requests),
		Errors:   atomic.LoadInt64(&s.Errors),
		Ready:    atomic.LoadInt32(&s.Ready),
	}
	return json.Marshal(snap)
}
