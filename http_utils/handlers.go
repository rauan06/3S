package http_utils

import (
	"net/http"
	"regexp"

	. "triples/bucket_struct"
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
	URL := r.URL.Path[1:]

	if SessionUser == nil {
		SessionUser = NewUser()
	}

	switch method {
	case "PUT":
		PUT(w, r, URL)
		return

	case "GET":
		GET(w, r, URL)
		return

	case "DELETE":
		DELETE(w, r, URL)
		return

	default:
		MethodNotAllowed(w, r)
		return
	}
}
