package commons

import (
	"log"
	"os"
	"path/filepath"
)

const (
	SLV            = "SLV"
	Version        = "1"
	appName        = "slv"
	appDataPathEnv = "SLV_APP_DATA"
)

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

var appDataDir *string

func AppDataDir() string {
	if appDataDir == nil {
		slvAppDataRoot := getAppDataDirPath()
		appDataDir = &slvAppDataRoot
	}
	return *appDataDir
}
