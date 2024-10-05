package http_utils

import (
	"encoding/xml"
	"fmt"
	"net/http"

	. "triples/buckets"
)

func PUT(w http.ResponseWriter, r *http.Request, bucketName string) {
	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}
	if !BucketIsUnique(bucketName) {
		ConflictRequest(w, r)
		return
	}

	// TODO: Add data reading
	newBucket := NewBucket(bucketName, nil)
	BucketNames = append(BucketNames, bucketName)

	AllBuckets[SessionUser.UserID] = append(AllBuckets[SessionUser.UserID], newBucket)

	OkRequestWithHeaders(w, r)
}

func GET(w http.ResponseWriter, r *http.Request, bucketName string) {
	if len(bucketName) != 0 {
		if !CheckRegex(bucketName) {
			BadRequest(w, r)
			return
		}
	}

	if len(BucketNames) == 0 {
		OkRequest(w, r)
		return
	}

	ListAllMyAllBucketsResult := NestForXML()

	OkRequestWithHeaders(w, r)
	out, _ := xml.MarshalIndent(ListAllMyAllBucketsResult, "  ", "  ")
	fmt.Fprint(w, xml.Header)
	fmt.Fprintln(w, string(out))
}

func DELETE(w http.ResponseWriter, r *http.Request, bucketName string) {
	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}
	for i, bucket := range AllBuckets[SessionUser.UserID] {
		if bucketName == bucket.Name {
			if len(bucket.Data) != 0 {
				ConflictRequest(w, r)
				return
			} else {
				DeleteFromBucketNames(bucket.Name)
				AllBuckets[SessionUser.UserID] = append(AllBuckets[SessionUser.UserID][:i], AllBuckets[SessionUser.UserID][i+1:]...)
				NoContentRequest(w, r)
				return
			}
		}
	}

	NotFoundRequest(w, r)
}
