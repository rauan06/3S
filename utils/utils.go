package utils

import (
	"crypto/md5"
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
