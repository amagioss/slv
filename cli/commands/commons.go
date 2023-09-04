package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func PrintErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, color.RedString("error:\t"+err.Error()))
	os.Exit(1)
}

func PrintErrorMessageAndExit(errMessage string) {
	fmt.Fprintln(os.Stderr, color.RedString("error:\t"+errMessage))
	os.Exit(1)
}

func getEnvSecretKey() string {
	secretKey := os.Getenv("SLV_SECRET_KEY")
	if secretKey == "" {
		PrintErrorMessageAndExit("SLV_SECRET_KEY environment variable is not set")
	}
	return secretKey
}
