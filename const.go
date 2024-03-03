package slv

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"savesecrets.org/slv/core/config"
	"savesecrets.org/slv/core/crypto"
)

func VersionInfo() string {
	if appInfo == nil {
		appInfo = new(string)
		buildAt := "unknown"
		if builtAtTime, err := time.Parse(time.RFC3339, BuildDate); err == nil {
			builtAtLocalTime := builtAtTime.Local()
			buildAt = builtAtLocalTime.Format("02 Jan 2006 03:04:05 PM MST")
		}
		appInfoBuilder := strings.Builder{}
		appInfoBuilder.WriteString(config.Art)
		appInfoBuilder.WriteString("\n")
		appInfoBuilder.WriteString(config.AppNameUpperCase + ": " + config.Description)
		appInfoBuilder.WriteString("\n")
		appInfoBuilder.WriteString("-------------------------------------------------")
		appInfoBuilder.WriteString("\n")
		appInfoBuilder.WriteString(fmt.Sprintf("SLV Version : %s\n", Version))
		appInfoBuilder.WriteString(fmt.Sprintf("Built At    : %s\n", buildAt))
		appInfoBuilder.WriteString(fmt.Sprintf("Git Commit  : %s\n", Commit))
		appInfoBuilder.WriteString(fmt.Sprintf("Web         : %s\n", config.Website))
		appInfoBuilder.WriteString(fmt.Sprintf("Platform    : %s\n", runtime.GOOS+"/"+runtime.GOARCH))
		appInfoBuilder.WriteString(fmt.Sprintf("Go Version  : %s", runtime.Version()))
		*appInfo = appInfoBuilder.String()
	}
	return *appInfo
}

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = ""
	secretKey *crypto.SecretKey
	appInfo   *string
)

// Errors
var (
	errEnvironmentAccessNotFound = errors.New("environment doesn't have access. please set the required environment variables")
)
