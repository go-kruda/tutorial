package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Property 1: go.work completeness
// For any Go module in the repo (identified by go.mod), go.work must contain a matching use directive.
// Validates: Requirements 1.2, 9.3
func TestProperty_GoWorkCompleteness(t *testing.T) {
	root := repoRoot(t)

	// Read go.work content
	goWorkPath := filepath.Join(root, "go.work")
	data, err := os.ReadFile(goWorkPath)
	if err != nil {
		t.Fatalf("failed to read go.work: %v", err)
	}
	goWorkContent := string(data)

	// Find all go.mod files in the repo (excluding tests/ and the root go.work.sum area)
	var moduleDirs []string
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip hidden directories, tests/, and vendor/
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == ".kiro" || name == "vendor" || name == "tests" {
				return filepath.SkipDir
			}
		}
		if d.Name() == "go.mod" {
			rel, err := filepath.Rel(root, filepath.Dir(path))
			if err != nil {
				return err
			}
			// Skip root-level go.mod if it exists
			if rel != "." {
				moduleDirs = append(moduleDirs, rel)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk repo: %v", err)
	}

	if len(moduleDirs) == 0 {
		t.Fatal("no Go modules found in the repository")
	}

	for _, dir := range moduleDirs {
		// go.work uses ./ prefix with forward slashes
		useDirective := "./" + filepath.ToSlash(dir)
		if !strings.Contains(goWorkContent, useDirective) {
			t.Errorf("go.work missing use directive for module at %s (expected %q)", dir, useDirective)
		}
	}

	t.Logf("verified %d Go modules are referenced in go.work", len(moduleDirs))
}
