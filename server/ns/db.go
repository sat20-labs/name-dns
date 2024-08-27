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
