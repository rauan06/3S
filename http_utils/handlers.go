package http_utils

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"

	. "triples/buckets"
)

var (
	validBucketNameRegex = regexp.MustCompile("^([a-z0-9.-]{3,63})$")
	ipAddressRegex       = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	doubleDashPeriod     = regexp.MustCompile(`[-]{2}|\.\.`)
	buckets              = make(map[int][]*Bucket)
	bucketNames          []string
	sessionUser          *User
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
		if !checkRegex(bucketName) {
			BadRequest(w, r)
			return
		}
		// Checking for uniqueness
		if !bucketIsUnique(bucketName) {
			ConflictRequest(w, r)
			return
		}

		// TODO: Add data reading
		newBucket := NewBucket(bucketName, nil)
		bucketNames = append(bucketNames, bucketName)

		buckets[sessionUser.UserID] = append(buckets[sessionUser.UserID], newBucket)

		OkRequestWithHeaders(w, r)
		return

	case "GET":
		if len(bucketNames) == 0 {
			OkRequest(w, r)
			return
		}

		ListAllMyBucketsResult := nestForXML()

		OkRequestWithHeaders(w, r)
		out, _ := xml.MarshalIndent(ListAllMyBucketsResult, "  ", "  ")
		fmt.Fprint(w, xml.Header)
		fmt.Fprintln(w, string(out))
		return

	case "DELETE":
		for _, bucket := range buckets[sessionUser.UserID] {
			if bucketName == bucket.Name {
				if len(bucket.Data) == 0 {
					ConflictRequest(w, r)
					return
				} else {
					bucket.Data = [][]byte{}
					NoContentRequest(w, r)
					return
				}
			}
		}

		NotFoundRequest(w, r)
		return
	}
}

func bucketIsUnique(bucketName string) bool {
	for _, name := range bucketNames {
		if name == bucketName {
			return false
		}
	}

	return true
}

func nestForXML() *ListAllMyBucketsResult {
	ListAllMyBucketsResult := &ListAllMyBucketsResult{}
	Buckets := &Buckets{Bucket: buckets[sessionUser.UserID]}
	ListAllMyBucketsResult.Buckets = *Buckets
	ListAllMyBucketsResult.User = *sessionUser

	return ListAllMyBucketsResult
}

func checkRegex(test string) bool {
	if ipAddressRegex.MatchString(test) {
		return false
	}

	if m := validBucketNameRegex.FindStringSubmatch(test); m == nil {
		return false
	}

	if doubleDashPeriod.MatchString(test) {
		return false
	}

	return true
}
