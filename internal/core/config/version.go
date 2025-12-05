package config

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func VersionInfo() string {
	if appInfo == nil {
		appInfo = new(string)
		var committedAt string
		if builtAtTime, err := time.Parse(time.RFC3339, commitDate); err == nil {
			builtAtLocalTime := builtAtTime.Local()
			committedAt = builtAtLocalTime.Format("02 Jan 2006 03:04:05 PM MST")
		}
		appInfoBuilder := strings.Builder{}
		appInfoBuilder.WriteString(ColorizedArt())
		appInfoBuilder.WriteString("\n")
		appInfoBuilder.WriteString(description)
		appInfoBuilder.WriteString("\n")
		appInfoBuilder.WriteString("-------------------------------------------------")
		appInfoBuilder.WriteString("\n")
		appInfoBuilder.WriteString(fmt.Sprintf("SLV Version  : %s\n", Version))
		appInfoBuilder.WriteString(fmt.Sprintf("Built At     : %s\n", committedAt))
		appInfoBuilder.WriteString(fmt.Sprintf("Release      : %s\n", releaseURL))
		appInfoBuilder.WriteString(fmt.Sprintf("Git Commit   : %s\n", fullCommit))
		appInfoBuilder.WriteString(fmt.Sprintf("Web          : %s\n", website))
		appInfoBuilder.WriteString(fmt.Sprintf("Platform     : %s\n", runtime.GOOS+"/"+runtime.GOARCH))
		appInfoBuilder.WriteString(fmt.Sprintf("Go Version   : %s", runtime.Version()))
		*appInfo = appInfoBuilder.String()
	}
	return *appInfo
}
