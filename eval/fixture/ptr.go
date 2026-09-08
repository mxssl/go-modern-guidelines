package legacy

import "time"

// Config holds optional client settings.
type Config struct {
	Timeout  *time.Duration
	Retries  *int32
	Insecure *bool
}

func durationPtr(d time.Duration) *time.Duration { return &d }
func int32Ptr(v int32) *int32                    { return &v }
func boolPtr(v bool) *bool                       { return &v }

// DefaultConfig returns the client defaults.
func DefaultConfig() Config {
	return Config{
		Timeout:  durationPtr(30 * time.Second),
		Retries:  int32Ptr(3),
		Insecure: boolPtr(false),
	}
}

// FastConfig returns settings for low-latency environments.
func FastConfig() Config {
	return Config{
		Timeout:  durationPtr(500 * time.Millisecond),
		Retries:  int32Ptr(1),
		Insecure: boolPtr(true),
	}
}
