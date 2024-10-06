package http_utils

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"triples/utils"

	. "triples/bucket_struct"
)

var (
	usernameRegex = regexp.MustCompile("^[a-zA-Z0-9._-]{3,20}$")
	passRegex     = regexp.MustCompile(`^[A-Za-z\d]{8,}$`)
)

func PUT(w http.ResponseWriter, r *http.Request, URL string, pathToDir string) {
	bucketName := URL

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

	if SessionUser == nil {
		cookieValue, err := r.Cookie("session_id")
		if err != nil {
			SessionUser = NewUser("cookie", utils.MdHashing("cookiepass"), pathToDir)
			value := SessionUser.UserID
			CookieID = value

			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: value,
			})
		} else {
			CookieID = cookieValue.Value
		}
		AllUsers = append(AllUsers, SessionUser)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: CookieID,
		})
	}

	newBucket := NewBucket(bucketName, SessionUser.UserID, nil, pathToDir)
	AllBuckets = append(AllBuckets, newBucket)

	if err := SaveBucketsToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	if err := SaveUsersToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	path := pathToDir + "/" + bucketName
	if err := os.Mkdir(path, 0o700); err != nil {
		InternalServerError(w, r)
		return
	}

	OkRequestWithHeaders(w, r)
	return
}

func GET(w http.ResponseWriter, r *http.Request,
	bucketName string,
) {
	// if !CheckRegex(bucketName) {
	// 	BadRequest(w, r)
	// 	return
	// }

	if len(AllBuckets) == 0 || SessionUser == nil {
		NoContentRequest(w, r)
		return
	}

	ListAllMyAllBucketsResult := NestForXML()

	OkRequestWithHeaders(w, r)
	out, _ := xml.MarshalIndent(
		ListAllMyAllBucketsResult,
		" ",
		"  ",
	)
	fmt.Fprint(w, xml.Header)
	fmt.Fprintln(w, string(out))
}

func DELETE(w http.ResponseWriter, r *http.Request,
	bucketName string,
) {
	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}
	for i, bucket := range AllBuckets {
		if bucketName == bucket.Name {
			if len(bucket.Data) != 0 {
				ConflictRequest(w, r)
				return
			} else {
				AllBuckets = append(
					AllBuckets[:i],
					AllBuckets[i+1:]...,
				)
				NoContentRequest(w, r)
				return
			}
		}
	}

	NotFoundRequest(w, r)
}

func POST(w http.ResponseWriter, r *http.Request,
	URL string, pathToDir string,
) {
	if SessionUser != nil {
		ImATeapotRequest(w, r)
		return
	}

	username := r.FormValue("username")
	pass := r.FormValue("passwords")

	if usernameRegex.MatchString(username) && passRegex.MatchString(pass) {
		BadRequest(w, r)
		return
	}

	for _, user := range AllUsers {
		if user.Username == username {
			if user.Password == utils.MdHashing(pass) {
				SessionUser = NewUser(username, pass, pathToDir)
				NoContentRequest(w, r)
				return
			} else {
				BadRequest(w, r)
				return
			}
		}
	}

	SessionUser = NewUser(username, pass, pathToDir)
	CookieID, _ = utils.GenerateToken(SessionUser.UserID)
	AllUsers = append(AllUsers, SessionUser)

	if err := SaveUsersToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	NoContentRequest(w, r)
	return
}
