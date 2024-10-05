package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

func MdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:])
}

func GenerateToken(userID string) (string, error) {
	randomBytes := make([]byte, 16) // 16 bytes = 128 bits of randomness

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err // Handle error
	}

	tokenData := userID + hex.EncodeToString(randomBytes)

	hash := sha256.New()
	hash.Write([]byte(tokenData))

	token := hash.Sum(nil)

	return hex.EncodeToString(token), nil
}

func Expiration() time.Time {
	return time.Now().Add(8766 * time.Hour)
}

func CheckForHelpAndExit() {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--h") || strings.HasPrefix(arg, "-h") {
			printHelp()
			os.Exit(0)
		}
	}
}

func printHelp() {
	fmt.Println(`Simple Storage Service.

**Usage:**
	triple-s [-port <N>] [-dir <S>]  
	triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`)
}
