package http_utils

import (
	"net/http"
	"regexp"

	. "triples/buckets"
)

// TODO: Update regex
var validBucketNameRegex = regexp.MustCompile("^([a-z0-9.-]{3,63})$")

var (
	buckets     = make(map[*User][]*Bucket)
	bucketNames []string
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	bucketName := r.URL.Path[1:]

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
		user := NewUser()

		buckets[user] = append(buckets[user], newBucket)

		OkRequest(w, r)
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
