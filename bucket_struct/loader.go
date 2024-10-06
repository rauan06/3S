package bucket_struct

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

func LoadIDs(pathToDir string) {
	idNode := &IDs{}

	IDs, err := os.ReadFile(pathToDir + "/id.xml")
	fmt.Println(pathToDir + "/id.xml")
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

	err = os.WriteFile("storage/id.xml", data, 0o600) // Use octal notation for permissions
	if err != nil {
		log.Printf("error writing to file: %v\n", err)
		return
	}
}

func fatal() {
	log.Fatal("Unable to load files")
	os.Exit(1)
}
