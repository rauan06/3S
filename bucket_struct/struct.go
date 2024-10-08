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
	Data         []*File   `xml:"Data"`
}

type File struct {
	Path string
}

type Bucket struct {
	Name         string    `xml:"BucketName"`
	PathToBucket string    `xml:"Path"`
	CreateDate   time.Time `xml:"CreationDate"`
	LastModified time.Time `xml:"LastModifiedDate"`
	LifeCycle    time.Time `xml:"ExpirationDate"`
	Status       string    `xml:"Status"`
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

func NewBucket(name string, userID string, data *File, pathToDir string) *SessionBucket {
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
		Data:         []*File{data},
	}
}

func NewUser(username, storageDir string, AllUsers []*User) *User {
	UserID++

	if username == "" {
		for NotUnique(AllUsers, username) {
			UserID++
		}
		username = utils.MdHashing(strconv.Itoa(UserID))
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

func NotUnique(AllUsers []*User, token string) bool {
	for _, user := range AllUsers {
		if user.Username == token {
			return true
		}
	}
	return false
}
