package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"triples/bucket_struct"
	"triples/utils"

	. "triples/http_utils"
)

func main() {
	utils.CheckForHelpAndExit()

	port := flag.Int("port", 8091, "Port number")
	dir := flag.String("dir", "buckets", "Path to the directory")

	flag.Parse()

	StorageDir = "storage/"
	syscall.Mkdir(StorageDir, 0o700)

	PathToDir = StorageDir + *dir
	if _, err := os.Stat(PathToDir); os.IsNotExist(err) {
		syscall.Mkdir(PathToDir, Mode)
	}

	LoadBuckets()
	bucket_struct.LoadIDs(StorageDir)

	http.HandleFunc("/", Handler)
	address := ":" + strconv.Itoa(*port)
	log.Printf("Starting development server at http://localhost%s/\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
