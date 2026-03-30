package tests

import (
	"strings"
	"testing"
)

// Property 9: Section README language content
// All section READMEs are written in English.
// This test verifies each section README contains substantive English content.
// Validates: Requirements 10.1 (updated: English instead of Thai)
func TestProperty_READMETHThaiContent(t *testing.T) {
	root := repoRoot(t)

	readmeFiles := collectREADMETHFiles(t, root)

	if len(readmeFiles) == 0 {
		t.Fatal("no section README files found")
	}

	for _, relPath := range readmeFiles {
		content := readFileContent(t, relPath)

		// Must contain substantive English content (at least some common words)
		hasEnglish := false
		for _, word := range []string{"section", "learn", "step", "run", "install", "create", "handler", "route"} {
			if strings.Contains(strings.ToLower(content), word) {
				hasEnglish = true
				break
			}
		}

		if !hasEnglish {
			t.Errorf("%s: does not contain substantive English content", relPath)
		}
	}

	t.Logf("verified %d section README files contain English content", len(readmeFiles))
}
