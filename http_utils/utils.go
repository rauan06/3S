package http_utils

import . "triples/bucket_struct"

func BucketIsUnique(bucketName string) bool {
	for _, name := range BucketNames {
		if name == bucketName {
			return false
		}
	}

	return true
}

func NestForXML() *ListAllMyBucketsResult {
	ListAllMyAllBucketsResult := &ListAllMyBucketsResult{}
	Buckets := &Buckets{Bucket: AllBuckets[SessionUser.UserID]}
	ListAllMyAllBucketsResult.Buckets = *Buckets
	ListAllMyAllBucketsResult.User = *SessionUser

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

func DeleteFromBucketNames(s string) {
	for i := range BucketNames {
		if BucketNames[i] == s {
			BucketNames = append(BucketNames[:i], BucketNames[i+1:]...)
		}
	}
}
