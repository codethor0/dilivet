package version

import "fmt"

// These are overridden at build time via -ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func String() string {
	return fmt.Sprintf("mldsa %s (commit %s) built %s", Version, Commit, Date)
}
