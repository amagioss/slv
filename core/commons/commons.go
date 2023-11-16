package commons

import (
	"os"
)

func FileExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && !f.IsDir()
}

func DirExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.IsDir()
}

func StringPtr(s string) *string {
	return &s
}
