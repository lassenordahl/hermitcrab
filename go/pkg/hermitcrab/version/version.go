package version

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// ParseVersion parses a version string into a semver.Version.
func ParseVersion(v string) (*semver.Version, error) {
	// Remove 'v' prefix if present.
	v = strings.TrimPrefix(v, "v")

	// Split the version string into parts.
	parts := strings.Split(v, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid version format: %s", v)
	}

	// Handle the main version part.
	mainParts := strings.Split(parts[0], ".")
	if len(mainParts) < 2 {
		mainParts = append(mainParts, "0") // Add .0 if not present
	}

	// Handle the UI version semantics.
	uiParts := strings.Split(parts[1], ".")
	if len(uiParts) != 2 || uiParts[0] != "ui" {
		return nil, fmt.Errorf("invalid UI version format: %s", parts[1])
	}

	// Convert UI version to an integer.
	uiVersion, err := strconv.Atoi(uiParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid UI version number: %s", uiParts[1])
	}

	// Construct the semver-compatible version string.
	// TODO(lasse): Don't introduce an empty patch release.
	semverStr := fmt.Sprintf("%s-ui.%d", strings.Join(mainParts, "."), uiVersion)

	return semver.NewVersion(semverStr)
}

// GetMajorVersion returns the major version string (e.g., "24.1").
func GetMajorVersion(v *semver.Version) string {
	return fmt.Sprintf("%d.%d", v.Major(), v.Minor())
}

// CompareVersions compares two versions and returns true if v1 is greater than v2.
func CompareVersions(v1, v2 *semver.Version) bool {
	return v1.GreaterThan(v2)
}
