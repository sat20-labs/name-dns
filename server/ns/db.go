package ns

import (
	"encoding/binary"

	"github.com/sat20-labs/name-dns/common"
	"go.etcd.io/bbolt"
)

const (
	BUCKET_NAME_COUNT           = "nameCount"
	BUCKET_COMMON_SUMMARY       = "commomSummary"
	KEY_TOTAL_NAME_ACCESS_COUNT = "totalNameAccessCount"
	KEY_TOTAL_NAME_COUNT        = "totalNameCount"
)

func (s *Service) setTotalNameCount(count uint64) error {
	countBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(countBytes, count)
	return common.PutBucket(s.DB, BUCKET_COMMON_SUMMARY, []byte(KEY_TOTAL_NAME_COUNT), countBytes)
}

func (s *Service) getTotalNameCount() (uint64, error) {
	value, err := common.GetBucket(s.DB, BUCKET_COMMON_SUMMARY, []byte(KEY_TOTAL_NAME_COUNT))
	if err != nil {
		return 0, err
	}
	count := uint64(0)
	if value != nil {
		count = binary.BigEndian.Uint64(value)
	}
	return count, err
}

func (s *Service) incTotalNameAccessCount() error {
	value, err := common.GetBucket(s.DB, BUCKET_COMMON_SUMMARY, []byte(KEY_TOTAL_NAME_ACCESS_COUNT))
	if err != nil {
		return err
	}
	count := uint64(1)
	if value != nil {
		count = binary.BigEndian.Uint64(value) + 1
	}
	countBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(countBytes, count)
	return common.PutBucket(s.DB, BUCKET_COMMON_SUMMARY, []byte(KEY_TOTAL_NAME_ACCESS_COUNT), countBytes)
}

func (s *Service) getTotalNameAccessCount() (uint64, error) {
	value, err := common.GetBucket(s.DB, BUCKET_COMMON_SUMMARY, []byte(KEY_TOTAL_NAME_ACCESS_COUNT))
	if err != nil {
		return 0, err
	}
	count := uint64(0)
	if value != nil {
		count = binary.BigEndian.Uint64(value) + 1
	}

	return count, err
}

func (s *Service) incNameCount(name string) error {
	value, err := common.GetBucket(s.DB, BUCKET_NAME_COUNT, []byte(name))
	if err != nil {
		return err
	}
	count := uint64(0)
	if value != nil {
		count = binary.BigEndian.Uint64(value) + 1
	}
	countBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(countBytes, count)
	return common.PutBucket(s.DB, BUCKET_NAME_COUNT, []byte(name), countBytes)
}

func (s *Service) getNameCounts(page, pageSize int) ([]NameCount, int, error) {
	var nameCounts []NameCount
	var total int

	err := s.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME_COUNT))
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
			countVal := binary.BigEndian.Uint64(v)
			nameCounts = append(nameCounts, NameCount{Name: name, Count: countVal})
			count++
		}

		total = bucket.Stats().KeyN
		return nil
	})

	return nameCounts, total, err
}
