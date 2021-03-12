package utils

import (
	"encoding/json"
	"log"

	badger "github.com/dgraph-io/badger"
)

//KnowNodes Struct
type KnowNodes struct {
	IP []string
}

// EncodeKnowNodes Function
func (ip KnowNodes) EncodeKnowNodes() []byte {
	data, err := json.Marshal(ip)
	if err != nil {
		panic(err)
	}

	return data
}

// DecodeKnowNodes Function
func DecodeKnowNodes(data []byte) (KnowNodes, error) {
	var i KnowNodes
	err := json.Unmarshal(data, &i)
	return i, err
}

//SetKnowNodes Function
func SetKnowNodes(data []string) bool {

	options, errOpt := Options()
	if errOpt != nil {
		//return errOpt
		log.Fatal(errOpt)
	}

	db, errDB := OpenDB(DbPath, options)
	if errDB != nil {
		//return errDB
		log.Fatal(errDB)
	}

	errUpd := db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte("known_nodes"), KnowNodes{
			IP: data,
		}.EncodeKnowNodes())
		errSet := txn.SetEntry(e)
		return errSet
	})
	if errUpd != nil {
		log.Fatal(errUpd)
		return false
	}
	return true
}

//GetKnowNodes Function
func GetKnowNodes(data string) KnowNodes {
	var valCopy []byte

	options, errOpt := Options()
	if errOpt != nil {
		//return errOpt
		log.Fatal(errOpt)
	}

	db, errDB := OpenDB(DbPath, options)
	if errDB != nil {
		//return errDB
		log.Fatal(errDB)
	}

	errView := db.View(func(txn *badger.Txn) error {
		item, errGet := txn.Get([]byte(data))
		if errGet != nil {
			return errGet
			//log.Fatal(errGet)
		}

		errValue := item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		if errValue != nil {
			log.Fatal(errValue)
		}

		valCopy, errValue = item.ValueCopy(nil)
		if errValue != nil {
			log.Fatal(errValue)
		}

		return nil
	})
	if errView != nil {
		//return errDB
		log.Fatal(errView)
	}

	result, errDec := DecodeKnowNodes(valCopy)
	if errDec != nil {
		log.Fatal(errDec)
	}
	return result
}
