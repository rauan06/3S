package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"triples/http_utils"
	"triples/utils"
)

var validBucketName = regexp.MustCompile("^([a-z0-9.-]{3,63})$")

func handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case "PUT":
		w.Header().Add("Connection", "close")
		w.Header().Add("Server", "triple-s")

		m := validBucketName.FindStringSubmatch(r.URL.Path[1:])
		if m == nil {
			http_utils.BadRequest(w, r)
			return
		}

		parentPath := "buckets"
		childPath := r.URL.Path[1:]
		fullPath := fmt.Sprintf("%s/%s", parentPath, childPath)

		err := utils.EnsureDirExists(fullPath)
		if err != nil {
			http_utils.ConflictRequest(w, r)
			return
		}
		w.Header().Add("Location", fullPath)
		http_utils.OkRequest(w, r)

		return
	case "GET":
		http_utils.OkRequest(w, r)
	case "DELETE":
		http_utils.OkRequest(w, r)
	default:
		http_utils.MethodNotAllowed(w, r)
	}
}

func main() {
	os.RemoveAll("buckets")
	http.HandleFunc("/", handler)
	address := ":8090"
	log.Printf("Starting development server at http://localhost%s/\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
