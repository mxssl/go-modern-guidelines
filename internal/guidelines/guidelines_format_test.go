package guidelines

import (
	"fmt"
	"go/format"
	"runtime"
	"strings"
	"testing"

	"github.com/JetBrains/go-modern-guidelines/internal/goversion"
)

// TestGuidelineExamplesAreFormattedGo checks that every example snippet parses as
// Go and is stored in gofmt form, so explain output is a usable style reference.
// Snippets are fragments, so each blank-line-separated block is formatted in
// whichever context parses: package-level declarations or function body statements.
func TestGuidelineExamplesAreFormattedGo(t *testing.T) {
	toolchainVersion, ok := toolchainLanguageVersion()
	if !ok {
		t.Skipf("cannot determine toolchain Go version from %q", runtime.Version())
	}

	skipped := 0
	for _, guideline := range modernGoGuidelines {
		if goversion.Compare(toolchainVersion, guideline.sinceVersion) < 0 {
			skipped += len(guideline.examples)
			continue
		}
		for i, example := range guideline.examples {
			for _, side := range []struct {
				name    string
				snippet string
			}{{"before", example.before}, {"after", example.after}} {
				formatted, err := formatGoSnippet(side.snippet)
				if err != nil {
					t.Errorf("guideline %q example %d %s does not parse as Go: %v\n%s",
						guideline.id, i+1, side.name, err, side.snippet)
					continue
				}
				if formatted != side.snippet {
					t.Errorf("guideline %q example %d %s is not gofmt-formatted\nhave:\n%s\nwant:\n%s",
						guideline.id, i+1, side.name, side.snippet, formatted)
				}
			}
		}
	}
	if skipped > 0 {
		t.Logf("skipped %d example snippets newer than the Go %s toolchain", skipped, toolchainVersion)
	}
}

// formatGoSnippet returns the gofmt form of an example snippet.
func formatGoSnippet(snippet string) (string, error) {
	blocks := strings.Split(snippet, "\n\n")
	formatted := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if out, err := formatGoDeclarations(block); err == nil {
			formatted = append(formatted, out)
			continue
		}
		out, err := formatGoStatements(block)
		if err != nil {
			return "", err
		}
		formatted = append(formatted, out)
	}
	return strings.Join(formatted, "\n\n"), nil
}

func formatGoDeclarations(block string) (string, error) {
	out, err := format.Source([]byte("package p\n" + block + "\n"))
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	lines = lines[1:]
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	return strings.Join(lines, "\n"), nil
}

func formatGoStatements(block string) (string, error) {
	out, err := format.Source([]byte("package p\nfunc _() {\n" + block + "\n}\n"))
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	start := -1
	for i, line := range lines {
		if line == "func _() {" {
			start = i + 1
			break
		}
	}
	if start < 0 || lines[len(lines)-1] != "}" {
		return "", fmt.Errorf("unexpected gofmt wrapper output:\n%s", out)
	}
	body := lines[start : len(lines)-1]
	for i, line := range body {
		body[i] = strings.TrimPrefix(line, "\t")
	}
	return strings.Join(body, "\n"), nil
}

// toolchainLanguageVersion reports the major.minor Go version of the running
// toolchain, so examples for a newer Go than the test runs on are not checked.
func toolchainLanguageVersion() (string, bool) {
	version := runtime.Version()
	index := strings.Index(version, "go1.")
	if index < 0 {
		return "", false
	}
	parts := strings.SplitN(version[index+len("go"):], ".", 3)
	if len(parts) < 2 {
		return "", false
	}
	minor := parts[1]
	for i, r := range minor {
		if r < '0' || r > '9' {
			minor = minor[:i]
			break
		}
	}
	candidate := parts[0] + "." + minor
	if !goversion.IsMajorMinor(candidate) {
		return "", false
	}
	return candidate, true
}
