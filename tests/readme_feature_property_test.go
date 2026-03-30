package tests

import (
	"strings"
	"testing"
)

// Property 6: Root README feature coverage
// For any feature keyword in {Wing Transport, Typed Handler, Auto CRUD, DI Container, SSE, MCP, OpenAPI, Middleware},
// root README must contain it (case-insensitive).
// Validates: Requirements 2.2
func TestProperty_RootREADMEFeatureCoverage(t *testing.T) {
	content := readFileContent(t, "README.md")
	lower := strings.ToLower(content)

	featureKeywords := []string{
		"Wing Transport",
		"Typed Handler",
		"Auto CRUD",
		"DI Container",
		"SSE",
		"MCP",
		"OpenAPI",
		"Middleware",
	}

	for _, keyword := range featureKeywords {
		if !strings.Contains(lower, strings.ToLower(keyword)) {
			t.Errorf("root README.md missing feature keyword %q (case-insensitive)", keyword)
		}
	}

	t.Logf("verified %d feature keywords present in root README", len(featureKeywords))
}
