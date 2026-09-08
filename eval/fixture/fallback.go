package legacy

import (
	"os"
	"sync/atomic"
)

// DiskReads counts how many times the on-disk config was consulted.
var DiskReads int64

// loadDefaultRegion reads the fallback region from disk. It is slow and is
// only expected to run when nothing else provides a region.
func loadDefaultRegion() string {
	atomic.AddInt64(&DiskReads, 1)
	b, err := os.ReadFile("/etc/legacy/region")
	if err != nil {
		return "us-east-1"
	}
	return string(b)
}

// ResolveRegion picks the region from the explicit flag, then the environment,
// then the on-disk default.
func ResolveRegion(flagRegion string) string {
	if flagRegion != "" {
		return flagRegion
	}
	if env := os.Getenv("LEGACY_REGION"); env != "" {
		return env
	}
	return loadDefaultRegion()
}
