package tests

import (
	"os"
	"strings"
	"testing"
)

// Property 7: Root README directory tree accuracy
// For any top-level section directory, root README must contain that directory name.
// Validates: Requirements 2.7
func TestProperty_RootREADMEDirectoryTreeAccuracy(t *testing.T) {
	root := repoRoot(t)
	content := readFileContent(t, "README.md")

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("failed to read repo root: %v", err)
	}

	var sectionDirs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		// Only check numbered section directories (00-*, 01-*, etc.)
		if len(name) >= 3 && name[0] >= '0' && name[0] <= '9' && name[1] >= '0' && name[1] <= '9' && strings.Contains(name, "-") {
			sectionDirs = append(sectionDirs, name)
		}
	}

	if len(sectionDirs) == 0 {
		t.Fatal("no section directories found in the repository root")
	}

	for _, dir := range sectionDirs {
		dirWithSlash := dir + "/"
		if !strings.Contains(content, dirWithSlash) && !strings.Contains(content, dir) {
			t.Errorf("root README.md missing directory %q in its directory tree", dir)
		}
	}

	t.Logf("verified %d section directories appear in root README", len(sectionDirs))
}
