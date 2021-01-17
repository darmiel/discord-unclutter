package database

import (
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	bolt "go.etcd.io/bbolt"
	"log"
)

//goland:noinspection GoUnhandledErrorResult
func BlocksGhostping(userID string, config *duconfig.Config) (block bool, err error) {
	db, err := open(config)
	if err != nil {
		return false, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, e := tx.CreateBucketIfNotExists(bucketName)
		if e != nil {
			err = e
			return e
		}

		get := bucket.Get([]byte(userID))
		if get == nil {
			block = false
			return nil
		}

		if len(get) == 1 {
			block = get[0] == 1
		} else {
			log.Println("WARN: Length of ghostping value for", userID, "not 1")
		}

		return nil
	})

	return
}
