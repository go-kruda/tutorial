package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Property 2: Coding section structure
// For any coding section directory (all sections except 00-why-kruda/), it must contain
// README.md, starter/, and complete/.
// Validates: Requirements 4.1, 5.1, 6.1, 7.2, 8.2
func TestProperty_CodingSectionStructure(t *testing.T) {
	root := repoRoot(t)

	// Collect all coding section directories.
	// Coding sections are:
	//   01-beginner, 02-auto-crud, 03-intermediate (top-level)
	//   04-advanced/* sub-sections
	//   05-production/* sub-sections
	var codingSections []string

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("failed to read repo root: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()

		// Skip non-section directories
		if !isNumberedSection(name) {
			continue
		}

		// 00-why-kruda is documentation-only, skip it
		if name == "00-why-kruda" {
			continue
		}

		// 04-advanced and 05-production have sub-sections
		if name == "04-advanced" || name == "05-production" {
			subEntries, err := os.ReadDir(filepath.Join(root, name))
			if err != nil {
				t.Errorf("failed to read %s: %v", name, err)
				continue
			}
			for _, sub := range subEntries {
				if sub.IsDir() {
					codingSections = append(codingSections, filepath.Join(name, sub.Name()))
				}
			}
		} else {
			codingSections = append(codingSections, name)
		}
	}

	if len(codingSections) == 0 {
		t.Fatal("no coding sections found in the repository")
	}

	for _, section := range codingSections {
		sectionPath := filepath.Join(root, section)

		// Must contain README.md
		readmePath := filepath.Join(sectionPath, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			t.Errorf("coding section %s missing README.md", section)
		}

		// Must contain starter/
		starterPath := filepath.Join(sectionPath, "starter")
		if info, err := os.Stat(starterPath); os.IsNotExist(err) || !info.IsDir() {
			t.Errorf("coding section %s missing starter/ directory", section)
		}

		// Must contain complete/
		completePath := filepath.Join(sectionPath, "complete")
		if info, err := os.Stat(completePath); os.IsNotExist(err) || !info.IsDir() {
			t.Errorf("coding section %s missing complete/ directory", section)
		}
	}

	t.Logf("verified %d coding sections have README.md, starter/, and complete/", len(codingSections))
}

// isNumberedSection checks if a directory name starts with a two-digit prefix (e.g., "01-").
func isNumberedSection(name string) bool {
	if len(name) < 3 {
		return false
	}
	return name[0] >= '0' && name[0] <= '9' &&
		name[1] >= '0' && name[1] <= '9' &&
		strings.Contains(name, "-")
}
