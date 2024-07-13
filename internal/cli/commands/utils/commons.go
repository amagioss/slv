package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func ExitOnError(err error) {
	fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
	ErroredExit()
}

func ExitOnErrorWithMessage(errMessage string) {
	fmt.Fprintln(os.Stderr, color.RedString(errMessage))
	ErroredExit()
}

func ErroredExit() {
	os.Exit(1)
}

func SafeExit() {
	os.Exit(0)
}
