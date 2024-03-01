package input

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func GetHiddenInput(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	input, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return input, err
}

func GetVisibleInput(prompt string) (string, error) {
	var input string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&input)
	return input, err
}

func GetConfirmation(prompt, allowFor string) (bool, error) {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	return strings.EqualFold(input, allowFor), err
}

func IsInteractive() error {
	if !term.IsTerminal(int(syscall.Stdin)) {
		return errNonInteractiveTerminal
	}
	return nil
}
