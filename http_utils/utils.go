package http_utils

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	. "triples/bucket_struct"
)

// func NestForXML() *ListAllMyBucketsResult {
// 	ListAllMyAllBucketsResult := &ListAllMyBucketsResult{}
// 	Buckets := &Buckets{Bucket: AllBuckets[SessionUser.UserID]}
// 	ListAllMyAllBucketsResult.Buckets = *Buckets
// 	ListAllMyAllBucketsResult.User = *SessionUser

// 	return ListAllMyAllBucketsResult
// }

func CheckRegex(test string) bool {
	if IpAddressRegex.MatchString(test) {
		return false
	}

	if m := ValidBucketNameRegex.FindStringSubmatch(test); m == nil {
		return false
	}

	if DoubleDashPeriod.MatchString(test) {
		return false
	}

	return true
}

func LoadBuckets() {
	users, err := os.ReadFile("buckets/users.xml")

	tempBuckets := &Buckets{}
	if err != nil {
	} else if len(users) != 0 {
		if err := xml.Unmarshal(users, tempBuckets); err != nil {
			Fatal()
		}
		AllBuckets = tempBuckets.List
	}

	buckets, err := os.ReadFile("buckets/buckets.xml")
	if err != nil {
	} else if len(buckets) != 0 {
		if err := xml.Unmarshal(buckets, &AllBuckets); err != nil {
			Fatal()
		}
	}
}

func Fatal() {
	log.Fatal("Unable to load files")
	os.Exit(1)
}

func SaveBucketsToXMLFile() error {
	if err := os.MkdirAll("buckets", 0o755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	output, err := xml.MarshalIndent(AllBuckets, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling XML: %w", err)
	}

	output = append([]byte(xml.Header), output...)

	err = os.WriteFile("buckets/buckets.xml", output, 0o644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	log.Printf("New XML saved to buckets.xml\n")
	return nil
}
