package httpclient

import (
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab/bucket"
)

type TestServerOpt func(*TestServerConfig)

type TestServerConfig struct {
	BucketManager bucket.BucketManager
	Logger        *log.Logger
	CacheDir      string
}

func WithBucketManager(bm bucket.BucketManager) TestServerOpt {
	return func(config *TestServerConfig) {
		config.BucketManager = bm
	}
}

func WithLogger(logger *log.Logger) TestServerOpt {
	return func(config *TestServerConfig) {
		config.Logger = logger
	}
}

func WithCacheDir(cacheDir string) TestServerOpt {
	return func(config *TestServerConfig) {
		config.CacheDir = cacheDir
	}
}

func NewTestServer(t *testing.T, opts ...TestServerOpt) *Server {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	config := &TestServerConfig{
		BucketManager: bucket.NewMockBucketManager(ctrl),
		Logger:        log.New(os.Stdout, "test: ", log.LstdFlags),
		CacheDir:      t.TempDir(),
	}

	for _, opt := range opts {
		opt(config)
	}

	server := NewServer(config.BucketManager, config.CacheDir, config.Logger)

	// Ensure routes are set up
	server.routes()

	return server
}
