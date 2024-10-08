package bucket_struct

import (
	"strconv"
	"strings"
	"time"
	"triples/utils"
)

type SessionBucket struct {
	BucketId     string    `xml:"BucketID"`
	SessionID    string    `xml:"SessionID"`
	Name         string    `xml:"BucketName"`
	PathToBucket string    `xml:"Path"`
	CreateDate   time.Time `xml:"CreationDate"`
	LastModified time.Time `xml:"LastModifiedDate"`
	LifeCycle    time.Time `xml:"ExpirationDate"`
	Status       string    `xml:"Status"`
	Data         [][]byte  `xml:"Data"`
}

type Bucket struct {
	Name         string    `xml:"BucketName"`
	PathToBucket string    `xml:"Path"`
	CreateDate   time.Time `xml:"CreationDate"`
	LastModified time.Time `xml:"LastModifiedDate"`
	LifeCycle    time.Time `xml:"ExpirationDate"`
	Status       string    `xml:"Status"`
	Data         [][]byte  `xml:"Data"`
}

type ListAllMyAllBucketsResult struct {
	Bucket []*Bucket
	Owner  *User
}

type User struct {
	UserID   int    `xml:"UserID"`
	Username string `xml:"Username"`
}

type Users struct {
	List []*User `xml:"User"`
}

type Buckets struct {
	List []*Bucket `xml:"Bucket"`
}

type SessionBuckets struct {
	List []*SessionBucket
}

type IDs struct {
	UserID   int
	BucketId int
}

type Response struct {
	Code    int
	Message string
}

var (
	BucketId = 0
	UserID   = 0
)

func NewBucket(name string, userID string, data [][]byte, pathToDir string) *SessionBucket {
	BucketId++

	hashedBucketId, _ := utils.GenerateToken(strconv.Itoa(BucketId))
	SaveIDs(storagePath(pathToDir))

	return &SessionBucket{
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

func NewUser(username, storageDir string) *User {
	UserID++

	if username == "" {
		username, _ = utils.GenerateToken(strconv.Itoa(UserID))
	}

	SaveIDs(storagePath(storageDir))

	return &User{
		UserID:   UserID,
		Username: username,
	}
}

func storagePath(path string) string {
	return strings.SplitAfterN(path, "/", 2)[0]
}
