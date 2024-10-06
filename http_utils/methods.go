package http_utils

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	. "triples/bucket_struct"
	"triples/utils"
)

var (
	usernameRegex = regexp.MustCompile("^[a-zA-Z0-9._-]{3,20}$")
	passRegex     = regexp.MustCompile(`^[A-Za-z\d]{8,}$`)
)

func PUT(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.SplitAfter(r.URL.Path[1:], "/")[0]

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

	if token := r.URL.Query().Get("session_id"); token != "" {
		for _, user := range AllUsers {
			if token == user.UserID {
				SessionUser = user
				break
			}
		}
	}

	if SessionUser == nil {
		SessionUser = NewUser("cookie", utils.MdHashing("cookiepass"), PathToDir)
	}

	newBucket := NewBucket(bucketName, SessionUser.UserID, nil, PathToDir)
	AllBuckets = append(AllBuckets, newBucket)

	if err := SaveBucketsToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	if err := SaveUsersToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	path := PathToDir + "/" + bucketName
	if err := os.Mkdir(path, 0o700); err != nil {
		ConflictRequest(w, r)
		return
	}

	sessioID := "Bucket session id: " + SessionUser.UserID
	respondWithXML(w, r, &Response{Code: 200, Messege: sessioID})
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

	if len(AllBuckets) == 0 || (SessionUser == nil && token == "") {
		NotFoundRequest(w, r)
		return
	}

	if len(bucketName) == 0 {
		result, _ := NestForXML(nil, nil)
		respondWithXML(w, r, result)
		return
	}

	switch len(pathParts) {
	case 1:
		handleBucketRequest(w, r, bucketName, token)
	case 2:
		fmt.Println(objectName)
		return
	}

	NotFoundRequest(w, r)
}

func handleBucketRequest(w http.ResponseWriter, r *http.Request, bucketName, token string) {
	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}

	var tempUser *User
	for _, user := range AllUsers {
		if user.UserID == token {
			tempUser = &User{}
			break
		}
	}

	for _, bucket := range AllBuckets {
		if bucket.Name == bucketName {
			result, _ := NestForXML(bucket, tempUser)
			respondWithXML(w, r, result)
			return
		}
	}

	NotFoundRequest(w, r)
}

func respondWithXML(w http.ResponseWriter, r *http.Request, result interface{}) {
	OkRequestWithHeaders(w, r)
	out, _ := xml.MarshalIndent(result, " ", "  ")
	fmt.Fprint(w, xml.Header)
	fmt.Fprintln(w, string(out))
}

func DELETE(w http.ResponseWriter, r *http.Request, bucketName string) {
	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}

	for i, bucket := range AllBuckets {
		if bucketName == bucket.Name {
			if len(bucket.Data) != 0 {
				ConflictRequest(w, r)
				return
			}
			AllBuckets = append(AllBuckets[:i], AllBuckets[i+1:]...)
			NoContentRequest(w, r)
			return
		}
	}

	NotFoundRequest(w, r)
}
