package http_utils

import (
	"encoding/xml"
	"log"
	"os"
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

	if err != nil {
		os.WriteFile("buckets/users.xml", nil, 666)
	} else if len(users) != 0 {
		if err := xml.Unmarshal(users, &AllUsers); err != nil {
			Fatal()
		}
	}

	buckets, err := os.ReadFile("buckets/buckets.xml")
	if err != nil {
		os.WriteFile("buckets/buckets.xml", nil, 666)
	} else if len(buckets) != 0 {
		if err := xml.Unmarshal(buckets, &AllUsers); err != nil {
			Fatal()
		}
	}
}

func Fatal() {
	log.Fatal("Unable to load files")
	os.Exit(1)
}
