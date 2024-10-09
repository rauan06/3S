package http_utils

import (
	"fmt"

	. "triples/bucket_struct"
)

func Logout() {
	SessionUser = nil
	AllBuckets = nil
	AllUsers = nil
}

func Login(token string) error {
	if SessionUser != nil && token == "" {
		return nil
	}

	if len(token) < 3 {
		return fmt.Errorf("Invalid token, token should be at least 3 characters long")
	}

	if len(token) > 64 {
		return fmt.Errorf("Invalid token, token cannot be more than 64 characters long")
	}

	if !CheckRegexToken(token) {
		return fmt.Errorf("Token may have only lowercase letters, numbers, hyphens, and periods")
	}

	for _, user := range AllUsers {
		if user.Username == token {
			SessionUser = user
			return nil
		}
	}

	tempUser := NewUser(token, StorageDir, AllUsers)
	AllUsers = append(AllUsers, tempUser)
	SessionUser = tempUser

	if err := SaveUsersToXMLFile(); err != nil {
		return fmt.Errorf("Srrver error, cant save users")
	}

	return nil
}
