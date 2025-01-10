// Package clip copies text to the system clipboard by shelling out to whatever
// clipboard utility the host provides.
package clip

import (
	"os/exec"
	"runtime"
)

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
