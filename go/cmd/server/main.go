package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab"
	"github.com/lassenordahl/hermitcrab/pkg/hermitcrab/bucket"
)

func main() {
	logger := log.New(os.Stdout, "hermitcrab: ", log.LstdFlags)

	testBucket := bucket.NewTestBucket()
	err := testBucket.AddVersion("v24.1-ui.01")
	if err != nil {
		logger.Fatalf("Failed to add version: %v", err)
	}

	crab, err := hermitcrab.New(hermitcrab.Config{
		BucketManager: testBucket,
		CacheDir:      "./cache",
		Logger:        logger,
	})
	if err != nil {
		logger.Fatalf("Failed to create HermitCrab: %v", err)
	}

	logger.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", crab.Server())
	if err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
