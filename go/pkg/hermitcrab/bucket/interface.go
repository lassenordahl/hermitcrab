package bucket

import (
	"context"
	"io"

	"github.com/Masterminds/semver/v3"
)

// BucketManager defines the interface for managing bucket operations.
type BucketManager interface {
	// GetLatestPatchVersion returns the latest patch version for a given major version.
	GetLatestPatchVersion(ctx context.Context, majorVersion string) (*semver.Version, error)

	// DownloadPatchVersion downloads the specified patch version and returns a reader.
	DownloadPatchVersion(ctx context.Context, version *semver.Version) (io.ReadCloser, error)
}
