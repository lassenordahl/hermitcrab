package version

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// ParseVersion parses a version string into a semver.Version
func ParseVersion(v string) (*semver.Version, error) {
	// Convert v24.1-ui.01 format to 24.1.0-ui.01 for semver compatibility
	parts := strings.Split(v, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid version format: %s", v)
	}

	version := strings.TrimPrefix(parts[0], "v")
	if !strings.Contains(version, ".") {
		version += ".0"
	}
	version += "-" + parts[1]

	return semver.NewVersion(version)
}

// GetMajorVersion returns the major version string (e.g., "24.1")
func GetMajorVersion(v *semver.Version) string {
	return fmt.Sprintf("%d.%d", v.Major(), v.Minor())
}

// CompareVersions compares two versions and returns true if v1 is greater than v2
func CompareVersions(v1, v2 *semver.Version) bool {
	return v1.GreaterThan(v2)
}
