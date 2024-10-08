package http_utils

import (
	"net/http"
	"regexp"

	. "triples/bucket_struct"
)

// Global variables
var (
	ValidBucketNameRegex = regexp.MustCompile("^([a-z0-9.-]{3,63})$")
	ValidTokenRegex      = regexp.MustCompile("^([a-z0-9.-]{3,64})$")
	IpAddressRegex       = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	DoubleDashPeriod     = regexp.MustCompile(`[-]{2}|\.\.`)
	AllBuckets           []*SessionBucket
	AllUsers             []*User
	SessionUser          *User
	PathToDir            string
	StorageDir           string
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case "PUT":
		PUT(w, r)
		return

	case "GET":
		GET(w, r)
		return

	case "DELETE":
		DELETE(w, r)
		return

	default:
		MethodNotAllowed(w, r)
		return
	}
}
