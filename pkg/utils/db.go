package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	badger "github.com/dgraph-io/badger"
)

//DbPath Const
const (
	//DbPath Const
	DbPath = "github.com/josercp/stone/badger/db"
)

//Options Func
func Options() (badger.Options, error) {
	options := badger.DefaultOptions(DbPath)
	options.ValueDir = DbPath
	options.Logger = nil
	return options, nil
}

//Retry Function
func Retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	return db, err
}

//OpenDB Function
func OpenDB(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := Retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	} else {
		return db, nil
	}
}
