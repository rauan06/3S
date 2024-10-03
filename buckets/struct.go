package buckets

import (
	"time"

	"triples/utils"
)

type Bucket struct {
	BucketId   int
	Name       string
	CreateDate time.Time
	Logs       []time.Time
	Data       []byte
}

type User struct {
	UserID   int
	Username string
	Password string
}

var (
	BucketId = 0
	UserID   = 0
)

func NewBucket(name string, data []byte) *Bucket {
	BucketId++

	return &Bucket{
		BucketId:   BucketId,
		Name:       name,
		CreateDate: time.Now(),
		Logs:       []time.Time{time.Now()},
		Data:       data,
	}
}

func NewUser() *User {
	UserID++

	return &User{
		UserID:   UserID,
		Username: "user",
		Password: utils.MdHashing("123123"),
	}
}
