package main

import (
	"bytes"
	"log"

	"github.com/cockroachdb/pebble"
	"github.com/linxGnu/grocksdb"
)

func main() {
	testComparer := pebble.DefaultComparer

	options := pebble.Options{
		Comparer: testComparer,
	}

	testMerger := pebble.DefaultMerger

	options.Merger = testMerger

	fromDB, err := pebble.Open("data/pebble", &options)
	if err != nil {
		log.Fatal(err)
	}
	defer fromDB.Close()

	mySecretKey := []byte("my secret key")
	mySecretValue := []byte("my secret value")

	err = fromDB.Set(mySecretKey, mySecretValue, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("key inserted")

	toDB, err := grocksdb.OpenDb(grocksdb.NewDefaultOptions(), "data/pebble")
	if err != nil {
		log.Fatal(err)
	}
	defer toDB.Close()

	value, err := toDB.Get(grocksdb.NewDefaultReadOptions(), mySecretKey)
	if err != nil {
		log.Fatal(err)
	}

	if bytes.Compare(value.Data(), mySecretValue) != 0 {
		log.Print("secret value does not match")
		log.Fatal(err)
	}

	log.Printf("secret value from rocksDB is: %s", string(value.Data()))
}
