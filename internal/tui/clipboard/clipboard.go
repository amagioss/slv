// Package clipboard provides a thin clipboard abstraction for the SLV TUI.
//
// It wraps github.com/atotto/clipboard for native OS clipboard access on
// desktops, and falls back to an OSC52 terminal escape sequence so that
// "copy" still reaches the user's local clipboard when SLV is run inside an
// SSH session, a container, or anywhere else without a usable OS clipboard.
//
// The exported API mirrors golang.design/x/clipboard for drop-in replacement:
// Init, Write, Read and the FmtText constant.
package clipboard

import (
	"encoding/base64"
	"fmt"
	"os"

	atotto "github.com/atotto/clipboard"
	"golang.org/x/term"
)

// Format is a placeholder kept for API compatibility with golang.design/x/clipboard.
// Only text is supported.
type Format int

const FmtText Format = 0

// Init is a no-op retained for API compatibility. The underlying backends are
// initialised lazily and have no global state to set up.
func Init() error { return nil }

// Write puts the given bytes onto the clipboard.
//
// Strategy:
//  1. Try the OS clipboard via atotto/clipboard (pbcopy / xclip / xsel /
//     wl-copy / Win32 syscalls).
//  2. If that fails (no clipboard tool, no $DISPLAY, headless container, SSH
//     session), emit an OSC52 escape sequence so the controlling terminal
//     writes to the user's local clipboard.
func Write(_ Format, data []byte) {
	text := string(data)
	if err := atotto.WriteAll(text); err == nil {
		return
	}
	emitOSC52(text)
}

// Read returns the clipboard contents. There is no portable way to read the
// user's local clipboard over OSC52 (most terminals disable read for security),
// so this only consults the OS clipboard. Returns nil if unavailable.
func Read(_ Format) []byte {
	text, err := atotto.ReadAll()
	if err != nil {
		return nil
	}
	return []byte(text)
}

// emitOSC52 writes an OSC52 "set clipboard" escape sequence to stderr if and
// only if stderr is attached to a terminal.
//
// The TTY guard keeps the operation completely silent in degenerate
// environments (redirected stderr, non-interactive containers, log capture)
// where the escape bytes would otherwise leak into a file or log stream.
//
// When running inside tmux, the sequence is wrapped in tmux's DCS passthrough
// so it reaches the outer terminal without relying on `set -g set-clipboard on`.
func emitOSC52(text string) {
	if !term.IsTerminal(int(os.Stderr.Fd())) {
		return
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	seq := fmt.Sprintf("\x1b]52;c;%s\x07", encoded)

	if os.Getenv("TMUX") != "" {
		seq = fmt.Sprintf("\x1bPtmux;\x1b%s\x1b\\", seq)
	}

	fmt.Fprint(os.Stderr, seq)
}
