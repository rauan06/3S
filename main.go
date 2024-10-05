package main

import (
	"log"
	"net/http"
	"os"

	"triples/bucket_struct"
	"triples/http_utils"
	"triples/utils"
)

func main() {
	utils.CheckForHelpAndExit()

	http_utils.LoadBuckets()
	bucket_struct.LoadIDs()

	http.HandleFunc("/", http_utils.Handler)
	os.Mkdir("buckets", 0o700)

	address := ":8091"
	log.Printf("Starting development server at http://localhost%s/\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
