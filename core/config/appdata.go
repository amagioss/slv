package config

import (
	"log"
	"os"
	"path/filepath"

	"savesecrets.org/slv/core/commons"
)

func GetAppDataDir() string {
	if appDataDir == nil {
		slvAppDataRoot := os.Getenv(envar_SLV_APP_DATA_DIR)
		var err error
		if slvAppDataRoot == "" {
			slvAppDataRoot, err = os.UserConfigDir()
			if err != nil {
				log.Fatalf("Error while getting slv app data path: %v", err)
			}
			slvAppDataRoot = filepath.Join(slvAppDataRoot, AppNameLowerCase)
		}
		if !commons.DirExists(slvAppDataRoot) {
			err := os.MkdirAll(slvAppDataRoot, 0755)
			if err != nil {
				log.Fatalf("Error in creating the app data directory: %v", err)
			}
		}
		appDataDir = &slvAppDataRoot
	}
	return *appDataDir
}

func ResetAppDataDir() error {
	return os.RemoveAll(GetAppDataDir())
}
