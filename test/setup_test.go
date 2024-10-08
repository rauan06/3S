package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"triples/bucket_struct"
	"triples/http_utils"

	. "triples/http_utils"
)

func SetupSuite(tb testing.TB) func(tb testing.TB) {
	StorageDir = "storage_test/"
	if err := os.Mkdir(StorageDir, 0o700); err != nil && !os.IsExist(err) {
		tb.Fatal(err)
	}

	PathToDir = StorageDir + "buckets"
	if _, err := os.Stat(PathToDir); os.IsNotExist(err) {
		if err := os.Mkdir(PathToDir, 0o700); err != nil {
			tb.Fatal(err)
		}
	}

	LoadBuckets()
	bucket_struct.LoadIDs(StorageDir)

	return func(tb testing.TB) {
		Logout()
		os.RemoveAll(StorageDir)
	}
}

func SetupWithSession(tb testing.TB) func(tb testing.TB) {
	StorageDir = "storage_test_session/"
	if err := os.Mkdir(StorageDir, 0o700); err != nil && !os.IsExist(err) {
		tb.Fatal(err)
	}

	PathToDir = StorageDir + "buckets"
	if _, err := os.Stat(PathToDir); os.IsNotExist(err) {
		if err := os.Mkdir(PathToDir, 0o700); err != nil {
			tb.Fatal(err)
		}
	}

	LoadBuckets()

	bucket_struct.LoadIDs(StorageDir)

	SessionUser = bucket_struct.NewUser("", StorageDir, AllUsers)

	requests := []*http.Request{
		httptest.NewRequest(http.MethodPut, "/123?session_id=rauan", nil),
		httptest.NewRequest(http.MethodPut, "/1234?session_id=123", nil),
		httptest.NewRequest(http.MethodPut, "/12345?session_id=123", nil),
		httptest.NewRequest(http.MethodPut, "/123456?session_id=123", nil),
		httptest.NewRequest(http.MethodPut, "/1234567?session_id=1234", nil),
	}

	for _, req := range requests {
		rr := httptest.NewRecorder()
		http_utils.Handler(rr, req)

	}

	return func(tb testing.TB) {
		StorageDir = "storage_test_session/"
		Logout()
		os.RemoveAll(StorageDir)
	}
}
