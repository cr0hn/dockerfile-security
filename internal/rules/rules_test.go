package rules

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadInternalAll(t *testing.T) {
	rules, err := LoadInternal("all")
	if err != nil {
		t.Fatalf("LoadInternal(all): %v", err)
	}
	if len(rules) != 35 {
		t.Errorf("expected 35 rules, got %d", len(rules))
	}
}

func TestLoadInternalCore(t *testing.T) {
	rules, err := LoadInternal("core")
	if err != nil {
		t.Fatalf("LoadInternal(core): %v", err)
	}
	if len(rules) != 10 {
		t.Errorf("expected 10 core rules, got %d", len(rules))
	}
	for _, r := range rules {
		if r.ID == "" || r.Description == "" || r.Regex == "" {
			t.Errorf("rule has empty fields: %+v", r)
		}
	}
}

func TestLoadInternalCredentials(t *testing.T) {
	rules, err := LoadInternal("credentials")
	if err != nil {
		t.Fatalf("LoadInternal(credentials): %v", err)
	}
	if len(rules) != 11 {
		t.Errorf("expected 11 credential rules, got %d", len(rules))
	}
}

func TestLoadInternalNone(t *testing.T) {
	rules, err := LoadInternal("none")
	if err != nil {
		t.Fatalf("LoadInternal(none): %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules for 'none', got %d", len(rules))
	}
}

func TestLoadInternalDefault(t *testing.T) {
	rules, err := LoadInternal("")
	if err != nil {
		t.Fatalf("LoadInternal(''): %v", err)
	}
	if len(rules) != 35 {
		t.Errorf("expected 35 rules for default, got %d", len(rules))
	}
}

func TestLoadInternalInvalid(t *testing.T) {
	_, err := LoadInternal("invalid")
	if err == nil {
		t.Error("expected error for invalid selection")
	}
}

func TestLoadInternalSecurity(t *testing.T) {
	rules, err := LoadInternal("security")
	if err != nil {
		t.Fatalf("LoadInternal(security): %v", err)
	}
	if len(rules) != 7 {
		t.Errorf("expected 7 security rules, got %d", len(rules))
	}
	for _, r := range rules {
		if r.ID == "" || r.Description == "" || r.Regex == "" {
			t.Errorf("rule has empty fields: %+v", r)
		}
	}
}

func TestLoadInternalMultiple(t *testing.T) {
	rules, err := LoadInternal("core,security")
	if err != nil {
		t.Fatalf("LoadInternal(core,security): %v", err)
	}
	if len(rules) != 17 {
		t.Errorf("expected 17 rules (10 core + 7 security), got %d", len(rules))
	}

	rules, err = LoadInternal("credentials,packages")
	if err != nil {
		t.Fatalf("LoadInternal(credentials,packages): %v", err)
	}
	if len(rules) != 15 {
		t.Errorf("expected 15 rules (11 credentials + 4 packages), got %d", len(rules))
	}

	rules, err = LoadInternal("security,packages,configuration")
	if err != nil {
		t.Fatalf("LoadInternal(security,packages,configuration): %v", err)
	}
	if len(rules) != 14 {
		t.Errorf("expected 14 rules (7+4+3), got %d", len(rules))
	}
}

func TestLoadFromFile(t *testing.T) {
	// Use the custom-rules.yaml from testdata
	path := filepath.Join("..", "..", "testdata", "custom-rules.yaml")
	rules, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 custom rule, got %d", len(rules))
	}
	if rules[0].ID != "custom-001" {
		t.Errorf("expected custom-001, got %s", rules[0].ID)
	}
}

func TestLoadFromFileMissing(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadExternalFile(t *testing.T) {
	// Create a temp file
	tmp, err := os.CreateTemp("", "rules-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	content := `- id: test-001
  description: Test rule
  regex: '(test)'
  reference: https://example.com
  severity: Low
`
	if _, err := tmp.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	rules, err := LoadExternal(tmp.Name())
	if err != nil {
		t.Fatalf("LoadExternal: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}

func TestRuleFieldsNotEmpty(t *testing.T) {
	allRules, err := LoadInternal("all")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range allRules {
		if r.ID == "" {
			t.Error("found rule with empty ID")
		}
		if r.Description == "" {
			t.Errorf("rule %s has empty description", r.ID)
		}
		if r.Regex == "" {
			t.Errorf("rule %s has empty regex", r.ID)
		}
		if r.Reference == "" {
			t.Errorf("rule %s has empty reference", r.ID)
		}
	}
}

func TestIssueFromRule(t *testing.T) {
	r := Rule{
		ID:          "test-001",
		Description: "Test",
		Regex:       "(foo)",
		Reference:   "https://example.com",
		Severity:    "High",
	}
	issue := IssueFromRule(r)
	if issue.ID != r.ID || issue.Description != r.Description || issue.Severity != r.Severity {
		t.Errorf("IssueFromRule mismatch: %+v", issue)
	}
}

func TestLoadFromURL(t *testing.T) {
	yamlContent := `- id: http-001
  description: Rule loaded from HTTP
  regex: '(test)'
  reference: https://example.com
  severity: Low
- id: http-002
  description: Second HTTP rule
  regex: '(foo)'
  reference: https://example.com
  severity: Medium
`
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(yamlContent))
	}))
	defer server.Close()

	// Test loading rules from the mock server
	rules, err := LoadFromURL(server.URL)
	if err != nil {
		t.Fatalf("LoadFromURL: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].ID != "http-001" {
		t.Errorf("expected http-001, got %s", rules[0].ID)
	}
	if rules[1].ID != "http-002" {
		t.Errorf("expected http-002, got %s", rules[1].ID)
	}
}

func TestLoadFromURLNotFound(t *testing.T) {
	// Create a mock HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := LoadFromURL(server.URL)
	if err == nil {
		t.Error("expected error for 404 response")
	}
}

func TestLoadFromURLInvalidYAML(t *testing.T) {
	// Create a mock HTTP server that returns invalid YAML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not: valid: yaml: ["))
	}))
	defer server.Close()

	_, err := LoadFromURL(server.URL)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadExternalWithURL(t *testing.T) {
	yamlContent := `- id: ext-001
  description: External rule
  regex: '(test)'
  reference: https://example.com
  severity: Low
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(yamlContent))
	}))
	defer server.Close()

	// LoadExternal should detect URL and use LoadFromURL
	rules, err := LoadExternal(server.URL)
	if err != nil {
		t.Fatalf("LoadExternal with URL: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].ID != "ext-001" {
		t.Errorf("expected ext-001, got %s", rules[0].ID)
	}
}
