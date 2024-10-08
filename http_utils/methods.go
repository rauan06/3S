package http_utils

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	. "triples/bucket_struct"
)

func PUT(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.SplitN(r.URL.Path[1:], "/", 2)
	bucketName := pathParts[0]
	var objectName string

	if len(pathParts) > 1 {
		objectName = pathParts[1]
	}

	token := r.URL.Query().Get("session_id")

	if !CheckRegex(bucketName) {
		BadRequest(w, r, "Incorrect bucket name")
		return
	}

	if SessionUser == nil && token == "" {
		SessionUser = NewUser("", StorageDir, AllUsers)
		AllUsers = append(AllUsers, SessionUser)
	} else {
		if err := Login(token); err != nil {
			writeXML(w, fmt.Sprintf("%v", err), 400)
			return
		}
	}

	switch len(pathParts) {
	case 1:
		handlePut(w, r, bucketName)
	case 2:
		handlePutObject(w, r, bucketName, objectName)
	}
}

func GET(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.SplitN(r.URL.Path[1:], "/", 2)
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
	case 2:
		// Handle object retrieval if needed
		fmt.Println(objectName)
	}
}

func DELETE(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.SplitN(r.URL.Path[1:], "/", 2)
	bucketName := pathParts[0]
	token := r.URL.Query().Get("session_id")

	if SessionUser == nil && token == "" {
		ForbiddenRequest(w, r)
		return
	}

	Login(token)

	if !CheckRegex(bucketName) {
		BadRequest(w, r, "Incorrect bucket name")
		return
	}

	for i, bucket := range AllBuckets {
		if bucket.Name == bucketName {
			if len(bucket.Data) != 0 {
				ConflictRequest(w, r, "Non-empty bucket")
				return
			}
			if bucket.SessionID != SessionUser.Username {
				ForbiddenRequestTokenInvalid(w, r)
				return
			}
			AllBuckets = append(AllBuckets[:i], AllBuckets[i+1:]...)

			if err := os.RemoveAll(bucket.PathToBucket); err != nil {
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

	NotFoundRequest(w, r)
}

func handleBucketRequest(w http.ResponseWriter, r *http.Request, bucketName string) {
	if !CheckRegex(bucketName) {
		BadRequest(w, r, "Incorrect bucket name")
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
	w.Header().Set("Content-Type", "application/xml")

	out, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal XML", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, xml.Header, string(out)+"\n")
}

func handlePut(w http.ResponseWriter, r *http.Request, bucketName string) {
	for _, bucket := range AllBuckets {
		if bucket.Name == bucketName {
			ConflictRequest(w, r, "Bucket name already exists")
			return
		}
	}

	newBucket := NewBucket(bucketName, SessionUser.Username, &File{}, PathToDir)
	AllBuckets = append(AllBuckets, newBucket)

	if err := os.Mkdir(PathToDir+"/"+bucketName, 0o700); err != nil {
		ConflictRequest(w, r, "Bucket name already exists")
		return
	}

	if err := SaveBucketsToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	respondWithXML(w, r, &Response{Code: 200, Message: fmt.Sprintf("Bucket session id: %s", SessionUser.Username)})
}

func handlePutObject(w http.ResponseWriter, r *http.Request, bucketName, objectName string) {
	for _, bucket := range AllBuckets {
		if bucket.Name == bucketName {
			if bucket.PathToBucket == "" {
				InternalServerError(w, r)
				return
			}

			contentType := r.Header.Get("Content-Type")
			extension, err := fileExtension(contentType)
			if err != nil {
				BadRequest(w, r, "Invalid file type")
				return
			}

			filePath := fmt.Sprintf("%s/%s%s", bucket.PathToBucket, objectName, extension)

			file, err := os.Create(filePath)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			defer file.Close()

			if _, err := io.Copy(file, r.Body); err != nil {
				BadRequest(w, r, "Error writing your file")
				return
			}

			existingPaths := make(map[string]struct{})
			for _, data := range bucket.Data {
				existingPaths[data.Path] = struct{}{}
			}

			if _, exists := existingPaths["/"+objectName+extension]; !exists {
				bucket.Data = append(bucket.Data, &File{Path: "/" + objectName + extension})
			}

			go func() {
				if err := SaveBucketsToXMLFile(); err != nil {
					fmt.Println("Error saving buckets:", err)
				}
			}()

			OkRequestWithHeaders(w, r)
			return
		}
	}
	NotFoundRequest(w, r)
}
