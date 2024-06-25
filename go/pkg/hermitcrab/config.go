package hermitcrab

import (
	"log"

	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab/bucket"
)

type Config struct {
	// BucketManager is the bucket manager to use.
	BucketManager bucket.BucketManager
	// CacheDir is the directory to cache patches in.
	CacheDir      string
	// Logger is the logger to use.
	Logger        *log.Logger
	// API prefix is the prefix for the version API.
	// This should follow the format `/{APIPrefix}/api/v1/`
	APIPrefix     string
}
