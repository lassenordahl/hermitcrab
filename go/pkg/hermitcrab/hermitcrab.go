package hermitcrab

import (
	"net/http"

	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab/httpclient"
)

type HermitCrab struct {
	config Config
	server *httpclient.Server
}

// New creates a new HermitCrab instance.
func New(config Config) (*HermitCrab, error) {
	server := httpclient.NewServer(config.BucketManager, config.CacheDir, config.Logger)

	// Set up routes
	return &HermitCrab{
		config: config,
		server: server,
	}, nil
}

func (hc *HermitCrab) Server() http.Handler {
	return hc.server
}
