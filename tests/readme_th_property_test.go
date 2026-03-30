package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Property 8: README_TH structural completeness
// For any README_TH file (section READMEs for sections 01-05), it must contain a time estimate
// and links to starter/ and complete/.
// README_TH files are: 01-beginner/README.md, 02-auto-crud/README.md, 03-intermediate/README.md,
// and all README.md files under 04-advanced/*/ and 05-production/*/.
// Validates: Requirements 4.5, 4.6, 10.3
func TestProperty_READMETHStructuralCompleteness(t *testing.T) {
	root := repoRoot(t)

	readmeTHFiles := collectREADMETHFiles(t, root)

	if len(readmeTHFiles) == 0 {
		t.Fatal("no README_TH files found")
	}

	// Pattern for time estimate: a number followed by a time-related word
	timePattern := regexp.MustCompile(`\d+\s*(นาที|minutes?|mins?|ชั่วโมง|hours?|hrs?)`)

	for _, relPath := range readmeTHFiles {
		content := readFileContent(t, relPath)

		// Must contain a time estimate
		if !timePattern.MatchString(strings.ToLower(content)) {
			t.Errorf("%s: missing time estimate (expected a number followed by a time unit)", relPath)
		}

		// Must contain a link to starter/
		if !strings.Contains(content, "starter/") && !strings.Contains(content, "starter)") {
			t.Errorf("%s: missing link to starter/", relPath)
		}

		// Must contain a link to complete/
		if !strings.Contains(content, "complete/") && !strings.Contains(content, "complete)") {
			t.Errorf("%s: missing link to complete/", relPath)
		}
	}

	t.Logf("verified %d README_TH files have time estimates and starter/complete links", len(readmeTHFiles))
}

// collectREADMETHFiles returns relative paths to all README_TH files (sections 01-05).
func collectREADMETHFiles(t *testing.T, root string) []string {
	t.Helper()

	var files []string

	// Top-level coding sections: 01-beginner, 02-auto-crud, 03-intermediate
	topLevel := []string{"01-beginner", "02-auto-crud", "03-intermediate"}
	for _, dir := range topLevel {
		readmePath := filepath.Join(dir, "README.md")
		fullPath := filepath.Join(root, readmePath)
		if _, err := os.Stat(fullPath); err == nil {
			files = append(files, readmePath)
		}
	}

	// Sub-sections under 04-advanced/
	advancedDir := filepath.Join(root, "04-advanced")
	if entries, err := os.ReadDir(advancedDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				readmePath := filepath.Join("04-advanced", e.Name(), "README.md")
				fullPath := filepath.Join(root, readmePath)
				if _, err := os.Stat(fullPath); err == nil {
					files = append(files, readmePath)
				}
			}
		}
	}

	// Sub-sections under 05-production/
	prodDir := filepath.Join(root, "05-production")
	if entries, err := os.ReadDir(prodDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				readmePath := filepath.Join("05-production", e.Name(), "README.md")
				fullPath := filepath.Join(root, readmePath)
				if _, err := os.Stat(fullPath); err == nil {
					files = append(files, readmePath)
				}
			}
		}
	}

	return files
}
