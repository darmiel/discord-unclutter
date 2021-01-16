package database

import (
	bolt "go.etcd.io/bbolt"
	"log"
	"time"
)

//goland:noinspection GoUnhandledErrorResult
func BlocksGhostping(userID string) (block bool, err error) {
	db, err := bolt.Open(
		Path,
		0666,
		&bolt.Options{Timeout: 1 * time.Second}, // if the db file has a lock, the application would hang until released
	)
	if err != nil {
		return false, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// check if bucket "ghostping-opt-out" exists
		bucket, e := tx.CreateBucketIfNotExists(OptOutBucket)
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
