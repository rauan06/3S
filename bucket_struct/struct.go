package bucket_struct

import (
	"strconv"
	"time"

	"triples/utils"
)

type Bucket struct {
	BucketId     string    `xml:"BucketID"`
	UserID       string    `xml:"-"`
	Name         string    `xml:"BucketName"`
	CreateDate   time.Time `xml:"CreationDate"`
	LastModified time.Time `xml:"LastModifiedDate"`
	LifeCycle    time.Time `xml:"ExpirationDate"`
	Status       string    `xml:"Status"`
	Data         [][]byte  `xml:"-"`
}

type ListAllMyAllBucketsResult struct {
	Buckets []*Bucket
	Owner   *User
}

type User struct {
	UserID   string `xml:"UserID"`
	Username string `xml:"Username"`
	Password string
}

type Users struct {
	List []*User `xml:"User"`
}

type Buckets struct {
	List []*Bucket `xml:"Bucket"`
}

type IDs struct {
	UserID   int
	BucketId int
}

type Response struct {
	Code    int
	Messege string
}

var (
	BucketId = 0
	UserID   = 0
)

func NewBucket(name string, userID string, data [][]byte) *Bucket {
	BucketId++

	hashedBucketId, _ := utils.GenerateToken(strconv.Itoa(BucketId))
	SaveIDs()

	return &Bucket{
		BucketId:     hashedBucketId,
		UserID:       userID,
		Name:         name,
		CreateDate:   time.Now(),
		LastModified: time.Now(),
		LifeCycle:    utils.Expiration(),
		Status:       "active",
		Data:         data,
	}
}

func NewUser(username, pass string) *User {
	UserID++

	SaveIDs()

	hashedUserId, _ := utils.GenerateToken(strconv.Itoa(UserID))

	return &User{
		UserID:   hashedUserId,
		Username: username,
		Password: utils.MdHashing(pass),
	}
}
