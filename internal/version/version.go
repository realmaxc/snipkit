// Package version exposes build and release metadata for snipkit.
//
// The exported variables are overridden at build time via -ldflags; see the
// Makefile and the release workflow for how they are populated.
package version

var (
	// Version is the semantic version of the build. It defaults to the last
	// tagged release for local `go build`/`go run` invocations.
	Version = "0.1.0"
	// Commit is the short git SHA the binary was built from.
	Commit = "none"
	// Date is the build timestamp in RFC 3339 form.
	Date = "unknown"
)

// String renders a single-line, human-readable version banner.
func String() string {
	return "snipkit " + Version + " (" + Commit + ", built " + Date + ")"
}
