package ns

import (
	"encoding/binary"

	"github.com/sat20-labs/name-ns/common"
	"go.etcd.io/bbolt"
)

const (
	BUCKET_NAME = "nameCounts"
)

func incrementNameCount(db *bbolt.DB, name string) error {
	value, err := common.GetBucket(db, BUCKET_NAME, []byte(name))
	if err != nil {
		return err
	}

	count := 1
	if value != nil {
		count = int(binary.BigEndian.Uint32(value)) + 1
	}
	countBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(countBytes, uint32(count))
	return common.PutBucket(db, BUCKET_NAME, []byte(name), countBytes)
}

func getNameCount(db *bbolt.DB, name string) (int, error) {
	value, err := common.GetBucket(db, BUCKET_NAME, []byte(name))
	if err != nil {
		return 0, err
	}

	count := 0
	if value != nil {
		count = int(binary.BigEndian.Uint32(value))
	}

	return count, err
}

func getNameCounts(db *bbolt.DB, page, pageSize int) ([]NameCount, int, error) {
	var nameCounts []NameCount
	var total int

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		cursor := bucket.Cursor()

		skip := (page - 1) * pageSize
		count := 0

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if skip > 0 {
				skip--
				continue
			}

			if count >= pageSize {
				break
			}

			name := string(k)
			countVal := int(binary.BigEndian.Uint32(v))
			nameCounts = append(nameCounts, NameCount{Name: name, Count: countVal})
			count++
		}

		total = bucket.Stats().KeyN
		return nil
	})

	return nameCounts, total, err
}
