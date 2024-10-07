package bucket_struct

import (
	"strconv"
	"strings"
	"time"
	"triples/utils"
)

type Bucket struct {
	BucketId     string    `xml:"BucketID"`
	SessionID    string    `xml:"SessionID"`
	Name         string    `xml:"BucketName"`
	PathToBucket string    `xml:"Path"`
	CreateDate   time.Time `xml:"CreationDate"`
	LastModified time.Time `xml:"LastModifiedDate"`
	LifeCycle    time.Time `xml:"ExpirationDate"`
	Status       string    `xml:"Status"`
	Data         [][]byte  `xml:"-"`
}

type ListAllMyAllBucketsResult struct {
	Bucket []*Bucket
	Owner  *User
}

type User struct {
	UserID   string `xml:"UserID"`
	Username string `xml:"Username"`
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

func NewBucket(name string, userID string, data [][]byte, pathToDir string) *Bucket {
	BucketId++

	hashedBucketId, _ := utils.GenerateToken(strconv.Itoa(BucketId))
	SaveIDs(storagePath(pathToDir))

	return &Bucket{
		BucketId:     hashedBucketId,
		SessionID:    userID,
		PathToBucket: pathToDir + "/" + name,
		Name:         name,
		CreateDate:   time.Now(),
		LastModified: time.Now(),
		LifeCycle:    utils.Expiration(),
		Status:       "active",
		Data:         data,
	}
}

func NewUser(username, pathToDir string) *User {
	UserID++

	SaveIDs(storagePath(pathToDir))

	hashedUserId, _ := utils.GenerateToken(strconv.Itoa(UserID))

	return &User{
		UserID:   hashedUserId,
		Username: username + strconv.Itoa(UserID),
	}
}

func storagePath(path string) string {
	return strings.SplitAfter(path, "/")[0]
}
