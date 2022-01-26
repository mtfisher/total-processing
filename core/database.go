package core

import (
	"log"

	"github.com/boltdb/bolt"
)

func (c Core) QueryDB(bucketName, key []byte) (val []byte) {
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		val = b.Get(key)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func (c Core) InsertDB(bucketName, key, value []byte) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}
