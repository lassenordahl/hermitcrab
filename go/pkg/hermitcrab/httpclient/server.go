package httpclient

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/gorilla/mux"
	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab/bucket"
)

type Server struct {
	router        *mux.Router
	bucketManager bucket.BucketManager
	cacheDir      string
	logger        *log.Logger
}

// NewServer creates a new HTTP server. It serves the latest version of the patch by default.
func NewServer(bm bucket.BucketManager, cacheDir string, logger *log.Logger) *Server {
	s := &Server{
		router:        mux.NewRouter(),
		bucketManager: bm,
		cacheDir:      cacheDir,
		logger:        logger,
	}
	s.routes()
	return s
}

// routes sets up the routes for the server.
func (s *Server) routes() {
	s.router.HandleFunc("/", s.serveLatestVersion).Methods("GET")
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// serveLatestVersion serves the latest version of the patch.
func (s *Server) serveLatestVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO(lasse): This is introducing a dependency on the bucket implementation working.
	// This should gracefully fail and use the most-recent-patch release.
	latestVersion, err := s.bucketManager.GetLatestPatchVersion(ctx, "24.1") // Hardcoded for demo
	if err != nil {
		http.Error(w, "Failed to get latest version", http.StatusInternalServerError)
		return
	}
	s.serveVersion(w, r, latestVersion)
}

// serveVersion serves the specified version of the patch.
func (s *Server) serveVersion(w http.ResponseWriter, r *http.Request, v *semver.Version) {
	ctx := r.Context()
	versionDir := filepath.Join(s.cacheDir, v.String())

	// Check if the version is already downloaded.
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		if err := s.downloadAndExtract(ctx, v, versionDir); err != nil {
			http.Error(w, "Failed to prepare version", http.StatusInternalServerError)
			return
		}
	}

	// Serve the version.
	http.FileServer(http.Dir(versionDir)).ServeHTTP(w, r)
}

// downloadAndExtract downloads and extracts the specified version. It stores the version in the cache directory.
// The version is stored in a directory named after the version string.
func (s *Server) downloadAndExtract(ctx context.Context, v *semver.Version, destDir string) error {
	reader, err := s.bucketManager.DownloadPatchVersion(ctx, v)
	if err != nil {
		return fmt.Errorf("failed to download version: %w", err)
	}
	defer reader.Close()

	// Create a gzip reader.
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	// Create the destination directory.
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract the tarball.
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar: %w", err)
		}

		target := filepath.Join(destDir, header.Name)

		// Ensure the target is within the destination directory.
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("failed to write file contents: %w", err)
			}
			f.Close()
		}
	}

	return nil
}
