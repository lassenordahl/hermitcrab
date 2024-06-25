package hermitcrab

import (
	"log"

	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab/bucket"
)

type Config struct {
	BucketManager bucket.BucketManager
	CacheDir      string
	Logger        *log.Logger
}
