package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func exitOnError(err error) {
	fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
	erroredExit()
}

func erroredExit() {
	os.Exit(1)
}

func safeExit() {
	os.Exit(0)
}
