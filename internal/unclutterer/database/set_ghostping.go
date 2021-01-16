package database

import (
	bolt "go.etcd.io/bbolt"
	"time"
)

//goland:noinspection GoUnhandledErrorResult
func SetBlocksGhostping(userID string, block bool) (err error) {
	db, err := bolt.Open(
		Path,
		0666,
		&bolt.Options{Timeout: 1 * time.Second}, // if the db file has a lock, the application would hang until released
	)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// check if bucket "ghostping-opt-out" exists
		bucket, e := tx.CreateBucketIfNotExists(OptOutBucket)
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
