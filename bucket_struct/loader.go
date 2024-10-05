package bucket_struct

import (
	"encoding/xml"
	"log"
	"os"
)

func LoadIDs() {
	idNode := &IDs{}

	IDs, err := os.ReadFile("buckets/id.xml")
	if err != nil {
		os.WriteFile("buckets/id.xml", nil, 666)
	} else if len(IDs) != 0 {
		err := xml.Unmarshal(IDs, &idNode)
		if err != nil {
			fatal()
		} else {
			UserID = idNode.UserID
			BucketId = idNode.BucketId
		}
	}
}

func fatal() {
	log.Fatal("Unable to load files")
	os.Exit(1)
}
