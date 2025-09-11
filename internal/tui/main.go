package tui

import (
	"fmt"
	"os"
)

// RunTUI starts the TUI application
func RunTUI() error {
	tui := NewTUI()
	return tui.Run()
}

// RunTUIWithErrorHandling runs TUI with error handling
func RunTUIWithErrorHandling() {
	if err := RunTUI(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI Error: %v\n", err)
		os.Exit(1)
	}
}
