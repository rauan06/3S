package http_utils

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

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

	if !CheckRegex(bucketName) {
		BadRequest(w, r, "Incorrect bucket name")
		return
	}

	switch len(pathParts) {
	case 1:
		handleBucketRequest(w, r, bucketName)
	case 2:
		handleGetObjects(w, r, bucketName, objectName)
	}
}

func DELETE(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.SplitN(r.URL.Path[1:], "/", 2)
	bucketName := pathParts[0]
	token := r.URL.Query().Get("session_id")
	var objectName string

	if len(pathParts) > 1 {
		objectName = pathParts[1]
	}

	if SessionUser == nil && token == "" {
		ForbiddenRequest(w, r)
		return
	}

	Login(token)

	if !CheckRegex(bucketName) {
		BadRequest(w, r, "Incorrect bucket name")
		return
	}

	switch len(pathParts) {
	case 1:
		for i, bucket := range AllBuckets {
			if bucket.Name == bucketName {
				if bucket.Data != nil && len(bucket.Data) != 0 && bucket.Data[0] != nil {
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
	case 2:
		for j, bucket := range AllBuckets {
			if bucket.Name == bucketName {
				if bucket.Data == nil {
					NotFoundRequest(w, r)
					return
				}

				for i, file := range bucket.Data {
					if file == nil {
						continue
					}
					if file.Name == objectName {
						if err := os.RemoveAll(bucket.PathToBucket + file.Path); err != nil {
							BadRequest(w, r, "File is not found or corrupted")
							return
						}

						bucket.Data = append(bucket.Data[:i], bucket.Data[i+1:]...)

						if bucket.Data == nil || (len(bucket.Data) == 1 && bucket.Data[0] == nil) {
							bucket.Data = nil
						}
						if err := SaveBucketsToXMLFile(); err != nil {
							InternalServerError(w, r)
							return
						}

						AllBuckets[j] = bucket

						NoContentRequest(w, r)
						return
					}
				}
			}
		}
	}

	NotFoundRequest(w, r)
}

func handleBucketRequest(w http.ResponseWriter, r *http.Request, bucketName string) {
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

	newBucket := NewBucket(bucketName, SessionUser.Username, nil, PathToDir)
	AllBuckets = append(AllBuckets, newBucket)

	if err := syscall.Mkdir(PathToDir+"/"+bucketName, Mode); err != nil {
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
	if !CheckRegex(objectName) {
		BadRequest(w, r, "Invalid object name")
		return
	}

	for _, bucket := range AllBuckets {
		if bucket == nil {
			continue // Skip nil buckets
		}

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

			fi, err := file.Stat()
			if err != nil {
				BadRequest(w, r, "Error finding your file")
				return
			}

			if bucket.Data == nil {
				bucket.Data = []*File{}
			}

			existingPaths := make(map[string]struct{})
			for _, data := range bucket.Data {
				if data != nil {
					existingPaths[data.Path] = struct{}{}
				}
			}

			if _, exists := existingPaths["/"+objectName+extension]; !exists {
				bucket.Data = append(bucket.Data, &File{Name: objectName, Path: "/" + objectName + extension, SizeInBytes: fi.Size()})
			}

			bucket.LastModified = time.Now()

			if err := SaveBucketsToXMLFile(); err != nil {
				fmt.Println("Error saving buckets:", err)
			}

			OkRequestWithHeaders(w, r)
			return
		}
	}
	NotFoundRequest(w, r)
}

func handleGetObjects(w http.ResponseWriter, r *http.Request, bucketName, objectName string) {
	for _, bucket := range AllBuckets {
		if bucket == nil {
			continue
		}

		if bucket.Name == bucketName {
			if bucket.Data == nil {
				continue
			}

			if len(bucket.Data) == 0 {
				continue
			}

			for _, path := range bucket.Data {
				if bucket.Data == nil {
					continue
				}
				if path.Name == objectName {
					filePath := bucket.PathToBucket + path.Path

					file, err := os.Open(filePath)
					if err != nil {
						http.Error(w, "File not found", http.StatusNotFound)
						return
					}
					defer file.Close()

					fileInfo, err := file.Stat()
					if err != nil {
						http.Error(w, "Could not retrieve file info", http.StatusBadRequest)
						return
					}

					w.Header().Set("Content-Type", getContentType(filePath))
					w.Header().Set("Location", filePath)

					http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
					return
				}
			}

			NotFoundRequest(w, r)
			return
		}
	}

	NotFoundRequest(w, r)
	return
}
