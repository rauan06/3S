package http_utils

import (
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
		// user := NewUser()

		buckets[1] = append(buckets[1], newBucket)

		OkRequest(w, r)
		return
	case "GET":
		response := []string{}

		for _, bucket := range buckets[1] {
			response = append(response, bucket.Name)
		}

		fmt.Fprintln(w, response)
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
