package input

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func HiddenInput(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	input, err := term.ReadPassword(syscall.Stdin)
	fmt.Println()
	return input, err
}

func VisibleInput(prompt string) (string, error) {
	var output string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&output)
	fmt.Println()
	return output, err
}

func IsAllowed() error {
	if !term.IsTerminal(syscall.Stdin) {
		return errNonInteractiveTerminal
	}
	return nil
}
