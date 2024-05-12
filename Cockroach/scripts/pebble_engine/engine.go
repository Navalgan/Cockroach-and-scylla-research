package pebble_engine

import (
	"main/scripts"

	"bytes"
	"log"
	"strconv"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/sstable"
	"github.com/cockroachdb/pebble/vfs"
)

var (
	ReplicasCount = 7
)

type InternalKV struct {
	Key   string
	Value string
}

func FindInternalKeyInSST(fileName string, key []byte, showAll bool) []InternalKV {
	fileOpen, err := vfs.Default.Open(fileName)

	readable, err := sstable.NewSimpleReadable(fileOpen)
	if err != nil {
		log.Fatal(err)
	}

	testComparer := pebble.DefaultComparer
	testComparer.Name = "cockroach_comparator"

	options := sstable.ReaderOptions{
		Comparer:   testComparer,
		MergerName: "cockroach_merge_operator",
		Merge:      pebble.DefaultMerger.Merge,
	}

	r, err := sstable.NewReader(readable, options)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	i, err := r.NewIter(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer i.Close()

	ikey, value := i.First()

	res := make([]InternalKV, 0)
	for ikey != nil {
		if bytes.Contains(value.InPlaceValue(), key) {
			res = append(res, InternalKV{Key: ikey.String(), Value: string(value.InPlaceValue())})
			if !showAll {
				return res
			}
		}
		ikey, value = i.Next()
	}

	return res
}

func FindInternalKeyInNode(node string, key []byte, showAll bool) []InternalKV {
	testComparer := pebble.DefaultComparer
	testComparer.Name = "cockroach_comparator"

	options := pebble.Options{
		Comparer: testComparer,
	}

	testMerger := pebble.DefaultMerger
	testMerger.Name = "cockroach_merge_operator"

	options.Merger = testMerger

	db, err := pebble.Open(node, &options)
	if err != nil {
		log.Fatal(err)
	}

	i, err := db.NewIter(nil)
	if err != nil {
		log.Fatal(err)
	}

	res := make([]InternalKV, 0)
	for i.Next() {
		val, err := i.ValueAndErr()
		if err != nil {
			log.Fatal(err)
		}

		if bytes.Contains(val, key) {
			res = append(res, InternalKV{Key: string(i.Key()), Value: string(val)})
			if !showAll {
				return res
			}
		}
	}

	return res
}

func EngineFindInternalKeyInNode(node string, key []byte, showAll bool) []InternalKV {
	files := scripts.GetSSTFromNode(node)

	res := make([]InternalKV, 0)
	for _, file := range files {
		res = append(res, FindInternalKeyInSST(file, key, showAll)...)

		if !showAll {
			return res
		}
	}

	return res
}

func CheckNodes(nodes []string, key []byte, showAll bool) map[string][]InternalKV {
	res := make(map[string][]InternalKV)

	for _, node := range nodes {
		nodeRes := EngineFindInternalKeyInNode(node, key, showAll)

		if len(nodeRes) != 0 {
			res[node] = nodeRes

			if !showAll {
				return res
			}
		}
	}

	return res
}

func CheckAllNodes(key []byte, showAll bool) map[string][]InternalKV {
	res := make(map[string][]InternalKV)

	for i := 1; i <= ReplicasCount; i++ {
		node := "cdb0" + strconv.Itoa(i)

		nodeRes := EngineFindInternalKeyInNode(node, key, showAll)

		if len(nodeRes) != 0 {
			res[node] = nodeRes

			if !showAll {
				return res
			}
		}
	}

	return res
}
