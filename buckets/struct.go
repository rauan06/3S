package buckets

import (
	"time"

	"triples/utils"
)

type Bucket struct {
	BucketId     int    `xml:"BucketID"`
	Name         string `xml:"BucketName"`
	CreateDate   time.Time
	LastModified time.Time
	LifeCycle    time.Time
	Status       string
	Data         [][]byte
}

type User struct {
	UserID   int
	Username string
	Password string
}

type ListAllMyBucketsResult struct {
	Buckets Buckets
	User    User
}

type Buckets struct {
	Bucket []*Bucket
}

type Response string

var (
	BucketId = 0
	UserID   = 0
)

func NewBucket(name string, data [][]byte) *Bucket {
	BucketId++

	return &Bucket{
		BucketId:     BucketId,
		Name:         name,
		CreateDate:   time.Now(),
		LastModified: time.Now(),
		LifeCycle:    utils.Expiration(),
		Status:       "active",
		Data:         data,
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
