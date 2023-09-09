package commons

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var appDataDir *string

func getAppDataDirPath() (slvAppDataRoot string) {
	slvAppDataRoot = os.Getenv(appDataPathEnv)
	if slvAppDataRoot == "" {
		appDataDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("Error while getting slv app data path: %v", err)
		}
		slvAppDataRoot = filepath.Join(appDataDir, appName)
	}
	if f, err := os.Stat(slvAppDataRoot); err == nil && f.IsDir() {
		return
	}
	err := os.MkdirAll(slvAppDataRoot, 0755)
	if err != nil {
		log.Fatalf("Error in creating the app data directory: %v", err)
	}
	return
}

func AppDataDir() string {
	if appDataDir == nil {
		slvAppDataRoot := getAppDataDirPath()
		appDataDir = &slvAppDataRoot
	}
	return *appDataDir
}

func FileExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && !f.IsDir()
}

func DirExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.IsDir()
}

func SearchStruct(s interface{}, query string) bool {
	query = strings.ToLower(query)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Struct {
			if SearchStruct(f.Interface(), query) {
				return true
			}
		} else {
			if strings.Contains(strings.ToLower(f.String()), query) {
				return true
			}
		}
	}
	return false
}
