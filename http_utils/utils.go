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
	users, err := os.ReadFile("buckets/buckets.xml")
	if err != nil {
		log.Printf("Error reading buckets.xml: %v", err)
		return
	}

	tempBuckets := &Buckets{}
	if len(users) != 0 {
		if err := xml.Unmarshal(users, tempBuckets); err != nil {
			Fatal(fmt.Errorf("error unmarshalling buckets.xml: %w", err))
		}
		AllBuckets = append(AllBuckets, tempBuckets.List...)
	}

	buckets, err := os.ReadFile("buckets/users.xml")
	if err != nil {
		log.Printf("Error reading users.xml: %v", err)
		return
	}

	if len(buckets) != 0 {
		if err := xml.Unmarshal(buckets, &AllUsers); err != nil {
			Fatal(fmt.Errorf("error unmarshalling users.xml: %w", err))
		}
	}
}

func Fatal(err error) {
	log.Fatal(err)
}

func SaveBucketsToXMLFile() error {
	tempBuckets := &Buckets{List: AllBuckets}

	output, err := xml.MarshalIndent(tempBuckets, "", "  ")
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
