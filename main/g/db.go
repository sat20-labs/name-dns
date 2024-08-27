package g

import (
	"time"

	mainCommon "github.com/sat20-labs/name-ns/main/common"
	"go.etcd.io/bbolt"
)

var store *bbolt.DB

const (
	DB_NAME = "name.db"
)

func InitDB() (err error) {
	store, err = bbolt.Open(
		mainCommon.YamlCfg.DB.Path+DB_NAME,
		0600,
		&bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return
}

func ReleaseDB() {
	store.Close()
}

func UpdateDB(callback func(tx *bbolt.Tx) error) error {
	return store.Update(callback)
}
