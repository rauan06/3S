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

func SaveIDs() {
	idNode := &IDs{UserID: UserID, BucketId: BucketId}
	data, err := xml.MarshalIndent(idNode, " ", " ")
	if err != nil {
		log.Printf("error marshalling XML: %v\n", err)
		return
	}

	data = append([]byte(xml.Header), data...)

	err = os.WriteFile("buckets/id.xml", data, 0o600) // Use octal notation for permissions
	if err != nil {
		log.Printf("error writing to file: %v\n", err)
		return
	}

	log.Println("IDs saved successfully to buckets/id.xml")
}

func fatal() {
	log.Fatal("Unable to load files")
	os.Exit(1)
}
