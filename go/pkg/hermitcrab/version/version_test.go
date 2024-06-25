package version

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"Valid version", "24.1-ui.1", "24.1.0-ui.1", false},
		{"Valid version with v prefix", "v24.1-ui.1", "24.1.0-ui.1", false},
		{"Valid version with patch", "24.1.5-ui.1", "24.1.5-ui.1", false},
		{"Valid version with single digit", "24-ui.1", "24.0.0-ui.1", false},
		{"Invalid format", "24.1.ui.1", "", true},
		{"Invalid UI version", "24.1-ui.a", "", true},
		{"Missing UI version", "24.1-ui", "", true},
		{"Extra segments", "24.1.2-ui.1.extra", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseVersion(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestGetMajorVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Major.Minor", "24.1.0", "24.1"},
		{"Major.Minor.Patch", "24.1.5", "24.1"},
		{"Major only", "24.0.0", "24.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := semver.NewVersion(tt.input)
			assert.NoError(t, err)
			result := GetMajorVersion(v)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected bool
	}{
		{"v1 > v2", "24.1.0-ui.2", "24.1.0-ui.1", true},
		{"v1 < v2", "24.1.0-ui.1", "24.1.0-ui.2", false},
		{"v1 == v2", "24.1.0-ui.1", "24.1.0-ui.1", false},
		{"Major version difference", "25.0.0-ui.1", "24.1.0-ui.2", true},
		{"Minor version difference", "24.2.0-ui.1", "24.1.0-ui.2", true},
		{"Patch version difference", "24.1.1-ui.1", "24.1.0-ui.2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := semver.NewVersion(tt.v1)
			assert.NoError(t, err)
			v2, err := semver.NewVersion(tt.v2)
			assert.NoError(t, err)
			result := CompareVersions(v1, v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}
