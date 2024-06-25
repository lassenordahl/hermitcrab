package bucket

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Masterminds/semver/v3"
)

type TestBucket struct {
	versions map[string]*semver.Version
}

func NewTestBucket() *TestBucket {
	return &TestBucket{
		versions: make(map[string]*semver.Version),
	}
}

// AddVersion adds a version to the bucket. If the version is greater than the existing version
// for the same major version, it will replace the existing version.
func (tb *TestBucket) AddVersion(version string) error {
	v, err := semver.NewVersion(version)
	if err != nil {
		return err
	}
	majorVersion := fmt.Sprintf("%d.%d", v.Major(), v.Minor())
	if existing, ok := tb.versions[majorVersion]; !ok || v.GreaterThan(existing) {
		tb.versions[majorVersion] = v
	}
	return nil
}

// GetLatestPatchVersion returns the latest patch version for a given major version.
// This is a simple implementation that returns the version that was added last.
func (tb *TestBucket) GetLatestPatchVersion(ctx context.Context, majorVersion string) (*semver.Version, error) {
	v, ok := tb.versions[majorVersion]
	if !ok {
		return nil, fmt.Errorf("no version found for major version %s", majorVersion)
	}
	return v, nil
}

// DownloadPatchVersion downloads the specified patch version and returns a reader. The content
// of the reader is a simple HTML file with the version number. This is essentially operating as
// a test utility to create an extremely simple way to generate test clients.
func (tb *TestBucket) DownloadPatchVersion(ctx context.Context, version *semver.Version) (io.ReadCloser, error) {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	// Create a simple index.html file.
	content := fmt.Sprintf("<html><body><h1>Version: %s</h1></body></html>", version.String())
	hdr := &tar.Header{
		Name:    "index.html",
		Mode:    0600,
		Size:    int64(len(content)),
		ModTime: time.Now(),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gzw.Close(); err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
