package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"triples/bucket_struct"
	. "triples/http_utils"
	"triples/utils"
)

func main() {
	utils.CheckForHelpAndExit()

	port := flag.Int("port", 8091, "'--port N' Port number")
	dir := flag.String("dir", "buckets", "'--dir S' Path to the directory")

	flag.Parse()

	storageDir := "storage"
	os.Mkdir(storageDir, 0o700)

	PathToDir = storageDir + "/" + *dir
	if _, err := os.Stat(PathToDir); os.IsNotExist(err) {
		os.Mkdir(PathToDir, 0o700)
	}

	http.HandleFunc("/", Handler)
	LoadBuckets(PathToDir)
	bucket_struct.LoadIDs(PathToDir)

	address := ":" + strconv.Itoa(*port)
	log.Printf("Starting development server at http://localhost%s/\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
