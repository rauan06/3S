package http_utils

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	. "triples/bucket_struct"
)

func NestForXML() *ListAllMyAllBucketsResult {
	ListAllMyAllBucketsResult := &ListAllMyAllBucketsResult{}
	tempBuckets := []*Bucket{}

	for _, bucket := range AllBuckets {
		if bucket.UserID == SessionUser.UserID {
			tempBuckets = append(tempBuckets, bucket)
		}
	}

	ListAllMyAllBucketsResult.Buckets = tempBuckets
	ListAllMyAllBucketsResult.Owner = SessionUser

	return ListAllMyAllBucketsResult
}

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

func LoadBuckets(pathToDir string) {
	buckets, err := os.ReadFile(pathToDir + "/buckets.xml")
	if err != nil {
		log.Printf("Error reading buckets.xml: %v", err)
		return
	}

	tempBuckets := &Buckets{}
	if len(buckets) != 0 {
		if err := xml.Unmarshal(buckets, tempBuckets); err != nil {
			log.Fatalf("Error unmarshalling buckets.xml: %v", err)
		}
		AllBuckets = append(AllBuckets, tempBuckets.List...)
	}

	users, err := os.ReadFile(pathToDir + "/users.xml")
	if err != nil {
		log.Printf("Error reading users.xml: %v", err)
		return
	}

	tempUsers := &Users{}
	if len(users) != 0 {
		if err := xml.Unmarshal(users, tempUsers); err != nil {
			log.Fatalf("Error unmarshalling users.xml: %v", err)
		}
		AllUsers = append(AllUsers, tempUsers.List...)
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

	err = os.WriteFile("storage/buckets.xml", output, 0o644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func SaveUsersToXMLFile() error {
	tempUsers := &Users{List: AllUsers}

	output, err := xml.MarshalIndent(tempUsers, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling XML: %w", err)
	}

	output = append([]byte(xml.Header), output...)

	err = os.WriteFile("storage/users.xml", output, 0o644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
