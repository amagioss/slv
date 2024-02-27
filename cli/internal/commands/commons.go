package commands

import (
	"fmt"
	"os"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/term"
)

func getHiddenInputFromUser(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	input, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return input, err
}

func exitOnError(err error) {
	fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
	erroredExit()
}

func exitOnErrorWithMessage(errMessage string) {
	fmt.Fprintln(os.Stderr, color.RedString(errMessage))
	erroredExit()
}

func erroredExit() {
	os.Exit(1)
}

func safeExit() {
	os.Exit(0)
}
