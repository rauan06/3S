package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"triples/bucket_struct"
	"triples/http_utils"
	"triples/utils"
)

func main() {
	utils.CheckForHelpAndExit()

	http_utils.LoadBuckets()
	bucket_struct.LoadIDs()

	port := flag.Int("port", 8091, "'--port N' Port number")
	dir := flag.String("dir", "buckets", "'--dir S' Path to the directory")

	flag.Parse()

	http.HandleFunc("/", http_utils.Handler)
	os.Mkdir(*dir, 0o700)
	address := ":" + strconv.Itoa(*port)
	log.Printf("Starting development server at http://localhost%s/\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
