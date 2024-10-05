package http_utils

import (
	"net/http"
	"regexp"

	. "triples/buckets"
)

// Global variables
var (
	ValidBucketNameRegex = regexp.MustCompile("^([a-z0-9.-]{3,63})$")
	IpAddressRegex       = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	DoubleDashPeriod     = regexp.MustCompile(`[-]{2}|\.\.`)
	AllBuckets           = make(map[int][]*Bucket)
	BucketNames          []string
	SessionUser          *User
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	bucketName := r.URL.Path[1:]

	if SessionUser == nil {
		SessionUser = NewUser()
	}

	switch method {
	case "PUT":
		PUT(w, r, bucketName)
		return

	case "GET":
		GET(w, r, bucketName)
		return

	case "DELETE":
		DELETE(w, r, bucketName)
		return
	}
}
