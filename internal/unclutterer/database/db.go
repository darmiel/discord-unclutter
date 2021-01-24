package database

import (
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	bolt "go.etcd.io/bbolt"
	"time"
)

var bucketName []byte

func open(config *duconfig.Config) (*bolt.DB, error) {
	if bucketName == nil || len(bucketName) == 0 {
		bucketName = []byte(config.DatabaseBucketName)
	}
	return bolt.Open(
		config.DatabasePath,
		0666,
		&bolt.Options{Timeout: 1 * time.Second}, // if the db file has a lock, the application would hang until released
	)
}
