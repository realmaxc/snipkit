package clip

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCopyUsesDetectedTool(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("shim assumes the linux tool-lookup order")
	}
	dir := t.TempDir()
	captured := filepath.Join(dir, "captured.txt")
	shim := filepath.Join(dir, "xclip")
	script := "#!/bin/sh\nexec cat > " + captured + "\n"
	if err := os.WriteFile(shim, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+":"+os.Getenv("PATH"))

	if !Available() {
		t.Fatal("expected a clipboard tool to be available via the shim")
	}
	if err := Copy("hello clip"); err != nil {
		t.Fatalf("Copy: %v", err)
	}
	b, err := os.ReadFile(captured)
	if err != nil {
		t.Fatalf("reading captured file: %v", err)
	}
	if string(b) != "hello clip" {
		t.Errorf("clipboard captured %q, want %q", b, "hello clip")
	}
}
