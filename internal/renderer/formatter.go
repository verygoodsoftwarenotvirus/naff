package renderer

import (
	"bytes"
	"fmt"
	"go/format"
	"os/exec"
	"strings"
)

// FormatGo runs goimports and gofmt on Go source code.
// If goimports is not available, falls back to gofmt only.
func FormatGo(src []byte) ([]byte, error) {
	// Try goimports first (handles both imports and formatting).
	result, err := runGoimports(src)
	if err == nil {
		return result, nil
	}

	// Fall back to standard library gofmt.
	formatted, err := format.Source(src)
	if err != nil {
		return nil, fmt.Errorf("formatting Go source: %w\n\nsource:\n%s", err, indicateLines(src))
	}

	return formatted, nil
}

func runGoimports(src []byte) ([]byte, error) {
	path, err := exec.LookPath("goimports")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(src)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("goimports: %w: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// indicateLines adds line numbers to source code for error messages.
func indicateLines(src []byte) string {
	lines := strings.Split(string(src), "\n")
	var b strings.Builder
	for i, line := range lines {
		fmt.Fprintf(&b, "%4d: %s\n", i+1, line)
	}
	return b.String()
}
