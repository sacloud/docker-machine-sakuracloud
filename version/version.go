package version

import "fmt"

var (
	// Version return version no
	Version = "0.0.9"
	// GitCommit return Git commit
	GitCommit = "HEAD"
)

// FullVersion return full version string
func FullVersion() string {
	return fmt.Sprintf("%s, build %s", Version, GitCommit)
}
