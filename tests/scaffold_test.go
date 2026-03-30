package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot returns the path to the repository root (parent of tests/).
func repoRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join("..")
}

// readFileContent is a helper that reads a file relative to the repo root.
func readFileContent(t *testing.T, relPath string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), relPath))
	if err != nil {
		t.Fatalf("failed to read %s: %v", relPath, err)
	}
	return string(data)
}

// Validates: Requirements 1.1
func TestRootREADMEExists(t *testing.T) {
	path := filepath.Join(repoRoot(t), "README.md")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("root README.md does not exist: %v", err)
	}
	if info.IsDir() {
		t.Fatal("README.md is a directory, expected a file")
	}
}

// Validates: Requirements 1.3
func TestGitignoreContainsGoArtefacts(t *testing.T) {
	content := readFileContent(t, ".gitignore")

	patterns := []string{"*.exe", "*.test", "*.out", "vendor/"}
	for _, p := range patterns {
		if !strings.Contains(content, p) {
			t.Errorf(".gitignore missing pattern %q", p)
		}
	}
}

// Validates: Requirements 1.4
func TestSixTopLevelSectionDirsExist(t *testing.T) {
	dirs := []string{
		"00-why-kruda",
		"01-beginner",
		"02-auto-crud",
		"03-intermediate",
		"04-advanced",
		"05-production",
	}
	for _, d := range dirs {
		path := filepath.Join(repoRoot(t), d)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("section directory %s does not exist: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s exists but is not a directory", d)
		}
	}
}

// Validates: Requirements 2.5
func TestRootREADMEContainsPrerequisites(t *testing.T) {
	content := readFileContent(t, "README.md")

	prerequisites := []string{"Go 1.25", "Docker", "Git"}
	for _, p := range prerequisites {
		if !strings.Contains(content, p) {
			t.Errorf("root README.md missing prerequisite %q", p)
		}
	}
}

// Validates: Requirements 2.6
func TestRootREADMEContainsQuickStart(t *testing.T) {
	content := readFileContent(t, "README.md")

	commands := []string{"git clone", "go run"}
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Errorf("root README.md missing quick-start command %q", cmd)
		}
	}
}

// Validates: Requirements 3.1, 3.2
func TestWhyKrudaREADMEContainsComparisonTable(t *testing.T) {
	content := readFileContent(t, "00-why-kruda/README.md")

	frameworks := []string{"Kruda", "Fiber", "Echo", "Chi"}
	for _, fw := range frameworks {
		if !strings.Contains(content, fw) {
			t.Errorf("00-why-kruda/README.md missing framework column %q", fw)
		}
	}

	// Check for table structure (pipe-delimited rows)
	if !strings.Contains(content, "|") {
		t.Error("00-why-kruda/README.md does not appear to contain a comparison table")
	}
}

// Validates: Requirements 3.3
func TestWhyKrudaREADMEContainsBenchmarkFigures(t *testing.T) {
	content := readFileContent(t, "00-why-kruda/README.md")

	lower := strings.ToLower(content)
	if !strings.Contains(lower, "benchmark") {
		t.Error("00-why-kruda/README.md missing benchmark section")
	}

	// Check for numeric figures (throughput/latency numbers)
	hasNumbers := false
	for _, word := range []string{"req/s", "ms", "latency", "throughput"} {
		if strings.Contains(lower, word) {
			hasNumbers = true
			break
		}
	}
	if !hasNumbers {
		t.Error("00-why-kruda/README.md missing benchmark figures (expected req/s, ms, latency, or throughput)")
	}
}

// Validates: Requirements 4.2
func TestBeginnerREADMEStates30Minutes(t *testing.T) {
	content := readFileContent(t, "01-beginner/README.md")
	if !strings.Contains(content, "30") {
		t.Error("01-beginner/README.md does not state 30 minutes")
	}
}

// Validates: Requirements 5.2
func TestAutoCRUDREADMEStates30Minutes(t *testing.T) {
	content := readFileContent(t, "02-auto-crud/README.md")
	if !strings.Contains(content, "30") {
		t.Error("02-auto-crud/README.md does not state 30 minutes")
	}
}

// Validates: Requirements 6.2
func TestIntermediateREADMEStates45Minutes(t *testing.T) {
	content := readFileContent(t, "03-intermediate/README.md")
	if !strings.Contains(content, "45") {
		t.Error("03-intermediate/README.md does not state 45 minutes")
	}
}

// Validates: Requirements 6.5
func TestDockerComposeDefinesDatabaseService(t *testing.T) {
	content := readFileContent(t, "03-intermediate/docker-compose.yml")

	lower := strings.ToLower(content)
	if !strings.Contains(lower, "postgres") {
		t.Error("03-intermediate/docker-compose.yml does not define a PostgreSQL database service")
	}
	if !strings.Contains(lower, "services") {
		t.Error("03-intermediate/docker-compose.yml missing 'services' key")
	}
}

// Validates: Requirements 7.1
func TestAdvancedContains8SubSections(t *testing.T) {
	expected := []string{
		"01-di-container",
		"02-auth-middleware",
		"03-openapi",
		"04-sse",
		"05-mcp-server",
		"06-websocket",
		"07-testing",
		"08-architecture",
	}

	advancedDir := filepath.Join(repoRoot(t), "04-advanced")
	entries, err := os.ReadDir(advancedDir)
	if err != nil {
		t.Fatalf("failed to read 04-advanced/: %v", err)
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}

	if len(dirs) != 8 {
		t.Errorf("04-advanced/ contains %d sub-section directories, expected 8: %v", len(dirs), dirs)
	}

	for _, name := range expected {
		found := false
		for _, d := range dirs {
			if d == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("04-advanced/ missing expected sub-section directory %q", name)
		}
	}
}

// Validates: Requirements 8.1
func TestProductionContains3SubSections(t *testing.T) {
	expected := []string{
		"01-monitoring",
		"02-docker-deploy",
		"03-benchmark",
	}

	prodDir := filepath.Join(repoRoot(t), "05-production")
	entries, err := os.ReadDir(prodDir)
	if err != nil {
		t.Fatalf("failed to read 05-production/: %v", err)
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}

	if len(dirs) != 3 {
		t.Errorf("05-production/ contains %d sub-section directories, expected 3: %v", len(dirs), dirs)
	}

	for _, name := range expected {
		found := false
		for _, d := range dirs {
			if d == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("05-production/ missing expected sub-section directory %q", name)
		}
	}
}

// Validates: Requirements 7.9
func TestTestingCompleteContainsTestFiles(t *testing.T) {
	dir := filepath.Join(repoRoot(t), "04-advanced", "07-testing", "complete")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read 04-advanced/07-testing/complete/: %v", err)
	}

	hasTestFile := false
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), "_test.go") {
			hasTestFile = true
			break
		}
	}
	if !hasTestFile {
		t.Error("04-advanced/07-testing/complete/ does not contain any _test.go files")
	}
}

// Validates: Requirements 7.10
func TestArchitectureCompleteContainsLayerDirs(t *testing.T) {
	base := filepath.Join(repoRoot(t), "04-advanced", "08-architecture", "complete")
	layerDirs := []string{"handler", "service", "repository"}

	for _, d := range layerDirs {
		path := filepath.Join(base, d)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("04-advanced/08-architecture/complete/%s does not exist: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("04-advanced/08-architecture/complete/%s is not a directory", d)
		}
	}
}

// Validates: Requirements 8.4
func TestDockerDeployDockerfileUsesMultiStageBuild(t *testing.T) {
	content := readFileContent(t, "05-production/02-docker-deploy/complete/Dockerfile")

	// Multi-stage builds use multiple FROM instructions
	fromCount := 0
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(strings.ToUpper(line))
		if strings.HasPrefix(trimmed, "FROM ") {
			fromCount++
		}
	}

	if fromCount < 2 {
		t.Errorf("Dockerfile has %d FROM instructions, expected at least 2 for multi-stage build", fromCount)
	}
}
