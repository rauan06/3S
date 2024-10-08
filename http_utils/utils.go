package http_utils

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	. "triples/bucket_struct"
)

// TODO: Remove SessionID and Data from bucket's struct
func NestForXML(bucket *SessionBucket) (*ListAllMyAllBucketsResult, error) {
	if SessionUser == nil {
		return nil, fmt.Errorf("Invalid token")
	}

	var bucketsToProcess []*Bucket
	if bucket != nil {
		bucketsToProcess = append(bucketsToProcess, &Bucket{
			Name:         bucket.Name,
			PathToBucket: bucket.PathToBucket,
			CreateDate:   bucket.CreateDate,
			LastModified: bucket.LastModified,
			LifeCycle:    bucket.LifeCycle,
			Status:       bucket.Status,
		})
	} else {
		for _, bucket := range AllBuckets {
			if bucket.SessionID == SessionUser.Username {
				bucketsToProcess = append(bucketsToProcess, &Bucket{
					Name:         bucket.Name,
					PathToBucket: bucket.PathToBucket,
					CreateDate:   bucket.CreateDate,
					LastModified: bucket.LastModified,
					LifeCycle:    bucket.LifeCycle,
					Status:       bucket.Status,
				})
			}
		}
	}

	result := &ListAllMyAllBucketsResult{
		Bucket: bucketsToProcess,
		Owner:  SessionUser,
	}

	return result, nil
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

func CheckRegexToken(test string) bool {
	if IpAddressRegex.MatchString(test) {
		return false
	}

	if m := ValidTokenRegex.FindStringSubmatch(test); m == nil {
		return false
	}

	if DoubleDashPeriod.MatchString(test) {
		return false
	}

	return true
}

func LoadBuckets() {
	buckets, err := os.ReadFile(StorageDir + "/buckets.xml")
	if err != nil {
		log.Printf("Error reading buckets.xml: %v", err)
		return
	}

	tempBuckets := &SessionBuckets{}
	if len(buckets) != 0 {
		if err := xml.Unmarshal(buckets, tempBuckets); err != nil {
			log.Fatalf("Error unmarshalling buckets.xml: %v", err)
		}
		AllBuckets = append(AllBuckets, tempBuckets.List...)
	}

	users, err := os.ReadFile(StorageDir + "/users.xml")
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

func SaveBucketsToXMLFile() error {
	tempBuckets := &SessionBuckets{List: AllBuckets}

	output, err := xml.MarshalIndent(tempBuckets, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling XML: %w", err)
	}

	output = append([]byte(xml.Header), output...)

	err = os.WriteFile(StorageDir+"buckets.xml", output, 0o644)
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

	err = os.WriteFile(StorageDir+"users.xml", output, 0o644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
