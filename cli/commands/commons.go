package commands

import (
	"fmt"
	"os"
)

func PrintErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, red+"error:\t"+err.Error()+reset)
	os.Exit(1)
}

func PrintErrorMessageAndExit(errMessage string) {
	fmt.Fprintln(os.Stderr, red+"error:\t"+errMessage+reset)
	os.Exit(1)
}
