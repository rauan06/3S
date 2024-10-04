package http_utils

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"

	. "triples/buckets"
)

// TODO: Update regex
var validBucketNameRegex = regexp.MustCompile("^([a-z0-9.-]{3,63})$")

var (
	buckets     = make(map[int][]*Bucket)
	bucketNames []string
	sessionUser *User
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	bucketName := r.URL.Path[1:]

	if sessionUser == nil {
		sessionUser = NewUser()
	}

	switch method {
	case "PUT":
		// Checking for validity
		if m := validBucketNameRegex.FindStringSubmatch(bucketName); m == nil {
			BadRequest(w, r)
			return
		}

		// Checking for uniqueness
		if m := bucketIsUnique(bucketName); m == false {
			ConflictRequest(w, r)
			return
		}

		// TODO: Add data reading
		newBucket := NewBucket(bucketName, nil)

		buckets[sessionUser.UserID] = append(buckets[sessionUser.UserID], newBucket)

		OkRequestWithHeaders(w, r)
		return
	case "GET":
		if len(bucketNames) == 0 {
			OkRequest(w, r)
			return
		}
		ListAllMyBucketsResult := &ListAllMyBucketsResult{}
		Buckets := &Buckets{Bucket: buckets[sessionUser.UserID]}
		ListAllMyBucketsResult.Buckets = *Buckets
		ListAllMyBucketsResult.User = *sessionUser

		OkRequestWithHeaders(w, r)
		out, _ := xml.MarshalIndent(ListAllMyBucketsResult, "  ", "  ")
		fmt.Fprint(w, xml.Header)
		fmt.Fprintln(w, string(out))
		return
	}
}

func bucketIsUnique(bucketName string) bool {
	for _, name := range bucketNames {
		if name == bucketName {
			return false
		}
	}

	bucketNames = append(bucketNames, bucketName)
	return true
}
