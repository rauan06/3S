package test

import (
	"os"
	"testing"

	"triples/bucket_struct"
	. "triples/http_utils"
)

func SetupSuite(tb testing.TB) func(tb testing.TB) {
	StorageDir = "storage_test/"
	os.Mkdir(StorageDir, 0o700)

	PathToDir = StorageDir + "/buckets"
	if _, err := os.Stat(PathToDir); os.IsNotExist(err) {
		os.Mkdir(PathToDir, 0o700)
	}

	LoadBuckets()
	bucket_struct.LoadIDs(StorageDir)

	return func(tb testing.TB) {
		os.RemoveAll(StorageDir)
	}
}
