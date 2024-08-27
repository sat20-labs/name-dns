package common

import "go.etcd.io/bbolt"

func InitBucket(db *bbolt.DB, name string) (err error) {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
}

func PutBucket(db *bbolt.DB, bucketName string, key []byte, value []byte) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Put(key, value)
	})
}

func GetBucket(db *bbolt.DB, bucketName string, key []byte) ([]byte, error) {
	var ret []byte
	db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		ret = bucket.Get(key)
		return nil
	})
	return ret, nil
}

func DelBucket(db *bbolt.DB, bucketName string, key []byte) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Delete(key)
	})
}
