package http_utils

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	. "triples/bucket_struct"
	"triples/utils"
)

var (
	usernameRegex = regexp.MustCompile("^[a-zA-Z0-9._-]{3,20}$")
	passRegex     = regexp.MustCompile(`^[A-Za-z\d]{8,}$`)
)

func PUT(w http.ResponseWriter, r *http.Request,
	URL string,
) {
	log.Println(AllBuckets)
	bucketName := URL

	if SessionUser == nil {
		cookieValue, err := r.Cookie("session_id")
		if err != nil {
			SessionUser = NewUser("cookie", utils.MdHashing("cookiepass"))
			value := SessionUser.UserID
			CookieID = value

			fmt.Fprintf(w, "CookieID: %s\n", value)
		} else {
			CookieID = cookieValue.Value
		}
	}

	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}

	for _, bucket := range AllBuckets {
		if bucket.Name == bucketName {
			BadRequest(w, r)
			return
		}
	}
	newBucket := NewBucket(bucketName, SessionUser.UserID, nil)
	AllBuckets = append(AllBuckets, newBucket)

	if err := SaveBucketsToXMLFile(); err != nil {
		InternalServerError(w, r)
		return
	}

	OkRequestWithHeaders(w, r)
}

func GET(w http.ResponseWriter, r *http.Request,
	bucketName string,
) {
	if len(bucketName) != 0 {
		if !CheckRegex(bucketName) {
			BadRequest(w, r)
			return
		}
	}

	if len(AllBuckets) == 0 {
		OkRequest(w, r)
		return
	}

	// ListAllMyAllBucketsResult := NestForXML()

	// OkRequestWithHeaders(w, r)
	// out, _ := xml.MarshalIndent(
	// 	ListAllMyAllBucketsResult,
	// 	"  ",
	// 	"  ",
	// )
	// fmt.Fprint(w, xml.Header)
	// fmt.Fprintln(w, string(out))
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
	URL string,
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

	for _, user := range AllUsers.List {
		if user.Username == username {
			if user.Password == utils.MdHashing(pass) {
				SessionUser = NewUser(username, pass)
				NoContentRequest(w, r)
				return
			} else {
				BadRequest(w, r)
				return
			}
		}
	}

	SessionUser = NewUser(username, pass)
	CookieID, _ = utils.GenerateToken(SessionUser.UserID)
	NoContentRequest(w, r)
	return
}
