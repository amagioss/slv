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
