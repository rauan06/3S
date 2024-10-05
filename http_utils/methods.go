package http_utils

import (
	"fmt"
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
	bucketName := URL[1:]

	// Check if SessionUser is nil and handle accordingly
	if SessionUser == nil {
		cookieValue, err := r.Cookie("session_id")
		if err != nil {
			// Create a new user if no session cookie is found
			newUser := NewUser("cookie", utils.MdHashing("cookiepass"))
			value := newUser.UserID
			CookieID = value // Store the new session ID

			fmt.Fprintf(w, "You now have your new cookieID,\n"+
				"you'll need to provide this ID \n"+
				"if you want to work on the bucket\n"+
				"again, by using this command:\n"+
				"*curl -b \"session_id=your-session-id link*\"\n")
			fmt.Fprintf(w, "\n CookieID: %s\n", value)
			return // Exit early to avoid further processing
		}

		// If we do have a cookie, extract its value
		CookieID = cookieValue.Value
	}

	// Validate the bucket name against the regex
	if !CheckRegex(bucketName) {
		BadRequest(w, r)
		return
	}

	// Check for bucket existence
	for _, bucket := range AllBuckets.List {
		if bucket.Name == bucketName {
			BadRequest(w, r)
			return
		}
	}

	// Create a new bucket
	newBucket := NewBucket(bucketName, SessionUser.UserID, nil)
	AllBuckets.List = append(AllBuckets.List, newBucket)

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

	if len(AllBuckets.List) == 0 {
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
	for i, bucket := range AllBuckets.List {
		if bucketName == bucket.Name {
			if len(bucket.Data) != 0 {
				ConflictRequest(w, r)
				return
			} else {
				AllBuckets.List = append(
					AllBuckets.List[:i],
					AllBuckets.List[i+1:]...,
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
