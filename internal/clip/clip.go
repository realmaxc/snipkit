// Package clip copies text to the system clipboard by shelling out to whatever
// clipboard utility the host provides.
package clip

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
)

// ErrNoClipboard indicates no supported clipboard utility was found on PATH.
var ErrNoClipboard = errors.New("no clipboard utility found on PATH")

// command returns the clipboard-copy command and arguments for the current
// host, or ok=false when no supported utility is found on PATH.
func command() (name string, args []string, ok bool) {
	switch runtime.GOOS {
	case "darwin":
		if p, err := exec.LookPath("pbcopy"); err == nil {
			return p, nil, true
		}
	case "windows":
		if p, err := exec.LookPath("clip"); err == nil {
			return p, nil, true
		}
	default:
		// Prefer Wayland, then fall back to the X11 helpers.
		if p, err := exec.LookPath("wl-copy"); err == nil {
			return p, nil, true
		}
		if p, err := exec.LookPath("xclip"); err == nil {
			return p, []string{"-selection", "clipboard"}, true
		}
		if p, err := exec.LookPath("xsel"); err == nil {
			return p, []string{"--clipboard", "--input"}, true
		}
	}
	return "", nil, false
}

// Available reports whether a clipboard utility was found on PATH.
func Available() bool {
	_, _, ok := command()
	return ok
}

// Copy writes text to the system clipboard.
func Copy(text string) error {
	return CopyFrom(strings.NewReader(text))
}

// CopyFrom streams r into the system clipboard utility's standard input.
func CopyFrom(r io.Reader) error {
	name, args, ok := command()
	if !ok {
		return ErrNoClipboard
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = r
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard copy via %s: %w", name, err)
	}
	return nil
}
