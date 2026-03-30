package tests

import (
	"os"
	"path/filepath"
	"testing"
	"unicode"
)

// Property 10: README emoji presence
// For any README file in the repo, it must contain at least one emoji character.
// Validates: Requirements 10.4, 10.5
func TestProperty_READMEEmojiPresence(t *testing.T) {
	root := repoRoot(t)

	var readmeFiles []string
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
		if d.Name() == "README.md" {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			readmeFiles = append(readmeFiles, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk repo: %v", err)
	}

	if len(readmeFiles) == 0 {
		t.Fatal("no README files found in the repository")
	}

	for _, relPath := range readmeFiles {
		content := readFileContent(t, relPath)

		if !containsEmoji(content) {
			t.Errorf("%s: does not contain any emoji characters", relPath)
		}
	}

	t.Logf("verified %d README files contain emoji characters", len(readmeFiles))
}

// containsEmoji checks if a string contains at least one emoji character.
// Checks common emoji Unicode ranges including symbols, emoticons, dingbats,
// and supplemental symbols.
func containsEmoji(s string) bool {
	for _, r := range s {
		if isEmoji(r) {
			return true
		}
	}
	return false
}

// isEmoji checks if a rune is an emoji character.
func isEmoji(r rune) bool {
	// Common emoji ranges
	switch {
	// Miscellaneous Symbols and Pictographs
	case r >= 0x1F300 && r <= 0x1F5FF:
		return true
	// Emoticons
	case r >= 0x1F600 && r <= 0x1F64F:
		return true
	// Transport and Map Symbols
	case r >= 0x1F680 && r <= 0x1F6FF:
		return true
	// Supplemental Symbols and Pictographs
	case r >= 0x1F900 && r <= 0x1F9FF:
		return true
	// Symbols and Pictographs Extended-A
	case r >= 0x1FA00 && r <= 0x1FA6F:
		return true
	// Symbols and Pictographs Extended-B
	case r >= 0x1FA70 && r <= 0x1FAFF:
		return true
	// Dingbats
	case r >= 0x2702 && r <= 0x27B0:
		return true
	// Miscellaneous Symbols
	case r >= 0x2600 && r <= 0x26FF:
		return true
	// Enclosed Alphanumeric Supplement (circled letters, etc.)
	case r >= 0x2460 && r <= 0x24FF:
		return true
	// Box-drawing characters used as emoji
	case r >= 0x2500 && r <= 0x257F:
		return false
	// Variation selectors and other marks
	case r == 0xFE0F || r == 0x200D:
		return false
	// Check if it's in the Symbol category and above basic ASCII
	case r > 0x2000 && unicode.IsSymbol(r):
		return true
	}
	return false
}
