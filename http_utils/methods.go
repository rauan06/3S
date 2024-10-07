package http_utils

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	. "triples/bucket_struct"
)

var (
	usernameRegex = regexp.MustCompile("^[a-zA-Z0-9._-]{3,20}$")
	passRegex     = regexp.MustCompile(`^[A-Za-z\d]{8,}$`)
)

func PUT(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.SplitAfter(r.URL.Path[1:], "/")[0]
	token := r.URL.Query().Get("session_id")

	if SessionUser == nil && token == "" {
		ForbiddenRequest(w, r)
		return
	}

	Login(token)
	if err := SaveUsersToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}

	for _, bucket := range AllBuckets {
		if bucket.Name == bucketName {
			ConflictRequest(w, r)
			return
		}
	}

	newBucket := NewBucket(bucketName, SessionUser.UserID, nil, PathToDir)
	AllBuckets = append(AllBuckets, newBucket)

	if err := SaveBucketsToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	path := PathToDir + "/" + bucketName
	if err := os.Mkdir(path, 0o700); err != nil {
		ConflictRequest(w, r)
		return
	}

	sessionID := "Bucket session id: " + SessionUser.UserID
	respondWithXML(w, r, &Response{Code: 200, Messege: sessionID})
	return
}

func GET(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.SplitAfter(r.URL.Path[1:], "/")
	bucketName := pathParts[0]
	var objectName string

	if len(pathParts) > 1 {
		objectName = pathParts[1]
	}

	token := r.URL.Query().Get("session_id")

	if SessionUser == nil && token == "" {
		ForbiddenRequest(w, r)
		return
	}

	Login(token)
	if err := SaveUsersToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	if len(bucketName) == 0 {
		result, err := NestForXML(nil)
		if err != nil {
			ForbiddenRequestTokenInvalid(w, r)
			return
		}

		respondWithXML(w, r, result)
		return
	}

	switch len(pathParts) {
	case 1:
		handleBucketRequest(w, r, bucketName)
		return
	case 2:
		fmt.Println(objectName)
		return
	}

	NotFoundRequest(w, r)
	return
}

func DELETE(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.SplitAfter(r.URL.Path[1:], "/")
	bucketName := pathParts[0]

	token := r.URL.Query().Get("session_id")

	if SessionUser == nil && token == "" {
		ForbiddenRequest(w, r)
		return
	}

	Login(token)
	if err := SaveUsersToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	// TODO: Object deletion
	// var objectName string

	// if len(pathParts) > 1 {
	// 	objectName = pathParts[1]
	// }

	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}

	// TODO: Remove files from directory using Removeall()
	for i, bucket := range AllBuckets {
		if bucketName == bucket.Name {
			if len(bucket.Data) != 0 {
				ConflictRequest(w, r)
				return
			}

			if bucket.SessionID != SessionUser.UserID {
				ForbiddenRequestTokenInvalid(w, r)
				return
			} else {
				AllBuckets = append(AllBuckets[:i], AllBuckets[i+1:]...)

				if err := os.Remove(bucket.PathToBucket); err != nil {
					InternalServerError(w, r)
					return
				}

				if err := SaveBucketsToXMLFile(); err != nil {
					InternalServerError(w, r)
					return
				}

				NoContentRequest(w, r)
				return
			}
		}
	}

	NotFoundRequest(w, r)
	return
}

func handleBucketRequest(w http.ResponseWriter, r *http.Request, bucketName string) {
	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}

	for _, bucket := range AllBuckets {
		if bucket.Name == bucketName {
			result, err := NestForXML(bucket)
			if err != nil {
				ForbiddenRequestTokenInvalid(w, r)
				return
			}
			respondWithXML(w, r, result)
			return
		}
	}

	NotFoundRequest(w, r)
	return
}

func respondWithXML(w http.ResponseWriter, r *http.Request, result interface{}) {
	OkRequestWithHeaders(w, r)
	out, _ := xml.MarshalIndent(result, " ", "  ")
	fmt.Fprint(w, xml.Header)
	fmt.Fprintln(w, string(out))
}
