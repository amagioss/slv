package commands

import (
	"fmt"
	"os"
)

func PrintErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, red+"error:\t"+err.Error()+reset)
	os.Exit(1)
}
