package utils

import (
	"fmt"
	"os"
)

func EnsureDirExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create a bucket: %w", err)
		}
	} else {
		return fmt.Errorf("bucket already exists: %s\n", dirPath)
	}
	return nil
}
