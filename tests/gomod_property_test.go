package tests

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Property 3: go.mod standards
// For any go.mod file, it must declare Go >= 1.25 and require github.com/go-kruda/kruda.
// Validates: Requirements 9.1, 9.2
func TestProperty_GoModStandards(t *testing.T) {
	root := repoRoot(t)

	var goModFiles []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == ".kiro" || name == "vendor" || name == "tests" {
				return filepath.SkipDir
			}
		}
		if d.Name() == "go.mod" {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			// Skip root go.work-level go.mod if present
			if rel != "go.mod" {
				goModFiles = append(goModFiles, rel)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk repo: %v", err)
	}

	if len(goModFiles) == 0 {
		t.Fatal("no go.mod files found in the repository")
	}

	for _, modFile := range goModFiles {
		data, err := os.ReadFile(filepath.Join(root, modFile))
		if err != nil {
			t.Errorf("failed to read %s: %v", modFile, err)
			continue
		}
		content := string(data)

		// Check Go version >= 1.25
		goVersion := parseGoVersion(content)
		if goVersion == "" {
			t.Errorf("%s: missing go version directive", modFile)
		} else if !isGoVersionAtLeast125(goVersion) {
			t.Errorf("%s: go version %s is less than 1.25", modFile, goVersion)
		}

		// Check kruda dependency
		if !strings.Contains(content, "github.com/go-kruda/kruda") {
			t.Errorf("%s: missing require for github.com/go-kruda/kruda", modFile)
		}
	}

	t.Logf("verified %d go.mod files meet standards", len(goModFiles))
}

// parseGoVersion extracts the Go version from a go.mod file content.
func parseGoVersion(content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}
	}
	return ""
}

// isGoVersionAtLeast125 checks if a Go version string is >= 1.25.
func isGoVersionAtLeast125(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	return major > 1 || (major == 1 && minor >= 25)
}
