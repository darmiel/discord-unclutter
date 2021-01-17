package database

import (
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	bolt "go.etcd.io/bbolt"
)

//goland:noinspection GoUnhandledErrorResult
func SetBlocksGhostping(userID string, block bool, config *duconfig.Config) (err error) {
	db, err := open(config)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// check if bucket "ghostping-opt-out" exists
		bucket, e := tx.CreateBucketIfNotExists(bucketName)
		if e != nil {
			err = e
			return e
		}

		var val byte
		if block {
			val = 1
		} else {
			val = 0
		}

		err = bucket.Put([]byte(userID), []byte{val})
		return nil
	})

	return
}
