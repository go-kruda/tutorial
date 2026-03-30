package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Property 5: Starter TODO comments
// For any starter Go module, at least one .go file must contain `// TODO:`.
// Validates: Requirements 11.3
func TestProperty_StarterTODOComments(t *testing.T) {
	root := repoRoot(t)

	var starterDirs []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == ".kiro" || name == "vendor" || name == "tests" {
				return filepath.SkipDir
			}
			// Check if this is a starter directory containing a go.mod
			if name == "starter" {
				goModPath := filepath.Join(path, "go.mod")
				if _, err := os.Stat(goModPath); err == nil {
					rel, err := filepath.Rel(root, path)
					if err == nil {
						starterDirs = append(starterDirs, rel)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk repo: %v", err)
	}

	if len(starterDirs) == 0 {
		t.Fatal("no starter Go modules found in the repository")
	}

	for _, dir := range starterDirs {
		hasTODO := false
		dirPath := filepath.Join(root, dir)

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			t.Errorf("failed to read %s: %v", dir, err)
			continue
		}

		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dirPath, e.Name()))
			if err != nil {
				continue
			}
			if strings.Contains(string(data), "// TODO:") {
				hasTODO = true
				break
			}
		}

		if !hasTODO {
			t.Errorf("starter module %s has no .go file containing '// TODO:'", dir)
		}
	}

	t.Logf("verified %d starter modules contain TODO comments", len(starterDirs))
}
