package input

import (
	"fmt"
	"os"
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

func GetMultiLineHiddenInput(prompt string) (input []byte, err error) {
	if prompt != "" {
		fmt.Println(prompt)
	}
	fmt.Println("Press enter/return twice after finishing your input to submit.")
	var line []byte
	emptyLines := 0
	for {
		if line, err = term.ReadPassword(int(os.Stdin.Fd())); err != nil {
			return nil, err
		}
		if len(line) == 0 {
			emptyLines++
			if emptyLines == 2 {
				break
			}
		} else {
			for i := 0; i < emptyLines; i++ {
				input = append(input, '\n')
			}
			if len(input) > 0 {
				input = append(input, '\n')
			}
			input = append(input, line...)
			emptyLines = 0
		}
	}
	return input, err
}

func GetVisibleInput(prompt string) (string, error) {
	var input string
	if prompt != "" {
		fmt.Print(prompt)
	}
	_, err := fmt.Scanln(&input)
	return input, err
}

func ReadBufferFromStdin(prompt string) ([]byte, error) {
	var input []byte
	buffer := make([]byte, 1024)
	if prompt != "" {
		fmt.Println(prompt)
	}
	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil || n == 0 {
			break
		}
		input = append(input, buffer[:n]...)
	}
	return input, nil
}

func GetConfirmation(prompt, allowFor string) (bool, error) {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	return strings.EqualFold(input, allowFor), err
}

func IsInteractive() bool {
	return term.IsTerminal(int(syscall.Stdin))
}
