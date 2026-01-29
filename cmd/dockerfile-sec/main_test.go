package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cr0hn/dockerfile-sec/internal/rules"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary for testing
	tmpDir, err := os.MkdirTemp("", "dockerfile-sec-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath = filepath.Join(tmpDir, "dockerfile-sec")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(string(out) + err.Error())
	}

	os.Exit(m.Run())
}

func runCLI(args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binaryPath, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		exitCode = 1
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func runCLIWithStdin(input string, args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = strings.NewReader(input)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		exitCode = 1
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func TestBasicFileAnalysis(t *testing.T) {
	stdout, stderr, exitCode := runCLI("../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	// Output should be JSON (since stdout is piped)
	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("expected valid JSON output: %v\nGot: %s", err, stdout)
	}

	// Should find some issues in the example Dockerfile
	if len(issues) == 0 {
		t.Error("expected to find issues in Dockerfile-example")
	}

	// Verify specific issues are found
	// core-002 now only detects ENV/ARG/LABEL with credentials, not RUN commands
	// cred-001 should detect the --password in RUN command
	foundCred001 := false
	for _, issue := range issues {
		if issue.ID == "cred-001" {
			foundCred001 = true
		}
	}
	if !foundCred001 {
		t.Error("expected to find cred-001 (generic credential)")
	}
}

func TestStdinInput(t *testing.T) {
	dockerfile := `FROM python:3.8
RUN pip install --password SECRET pkg
`
	stdout, stderr, exitCode := runCLIWithStdin(dockerfile)
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}

	if len(issues) == 0 {
		t.Error("expected to find issues from stdin")
	}
}

func TestExitCodeFlag(t *testing.T) {
	// With -E flag and issues found, should exit 1
	_, _, exitCode := runCLI("-E", "../../testdata/Dockerfile-example")
	if exitCode != 1 {
		t.Errorf("expected exit code 1 with -E and issues, got %d", exitCode)
	}
}

func TestExitCodeFlagNoIssues(t *testing.T) {
	// With -E flag and -R none (no rules), should exit 0
	_, stderr, exitCode := runCLI("-E", "-R", "none", "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0 with -E and no rules (no issues), got %d, stderr: %s", exitCode, stderr)
	}
}

func TestQuietMode(t *testing.T) {
	stdout, stderr, exitCode := runCLI("-q", "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	// In quiet mode, stdout should be empty
	if stdout != "" {
		t.Errorf("expected empty stdout in quiet mode, got: %s", stdout)
	}
}

func TestOutputToFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "output-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, stderr, exitCode := runCLI("-o", tmpFile.Name(), "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	// Read the output file
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	var issues []rules.Issue
	if err := json.Unmarshal(data, &issues); err != nil {
		t.Fatalf("expected valid JSON in output file: %v", err)
	}

	if len(issues) == 0 {
		t.Error("expected issues in output file")
	}
}

func TestRuleSetsCore(t *testing.T) {
	stdout, stderr, exitCode := runCLI("-R", "core", "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("expected valid JSON: %v\nGot: %s", err, stdout)
	}

	// Should only find core issues, not credential issues
	for _, issue := range issues {
		if strings.HasPrefix(issue.ID, "cred-") {
			t.Errorf("found credential issue %s when only core rules were loaded", issue.ID)
		}
	}
}

func TestRuleSetsCredentials(t *testing.T) {
	stdout, stderr, exitCode := runCLI("-R", "credentials", "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("expected valid JSON: %v\nGot: %s", err, stdout)
	}

	// Should only find credential issues, not core issues
	for _, issue := range issues {
		if strings.HasPrefix(issue.ID, "core-") {
			t.Errorf("found core issue %s when only credential rules were loaded", issue.ID)
		}
	}
}

func TestRuleSetsNone(t *testing.T) {
	stdout, stderr, exitCode := runCLI("-R", "none", "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("expected valid JSON: %v\nGot: %s", err, stdout)
	}

	// Should find no issues with no rules
	if len(issues) != 0 {
		t.Errorf("expected 0 issues with -R none, got %d", len(issues))
	}
}

func TestIgnoreRuleByID(t *testing.T) {
	// Use Dockerfile-security-issues which has ENV with password (triggers core-002)
	// First run without ignore to get baseline
	stdout1, _, _ := runCLI("../../testdata/Dockerfile-security-issues")
	var baseline []rules.Issue
	json.Unmarshal([]byte(stdout1), &baseline)

	// Then run with ignore
	stdout2, stderr, exitCode := runCLI("-i", "core-002", "../../testdata/Dockerfile-security-issues")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout2), &issues); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}

	// Should have fewer issues
	if len(issues) >= len(baseline) {
		t.Errorf("expected fewer issues with -i core-002, got %d vs %d", len(issues), len(baseline))
	}

	// Should not contain core-002
	for _, issue := range issues {
		if issue.ID == "core-002" {
			t.Error("found ignored rule core-002 in output")
		}
	}
}

func TestIgnoreFile(t *testing.T) {
	stdout, stderr, exitCode := runCLI("-F", "../../testdata/ignore-1", "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}

	// Should not contain core-001 or core-007 (from ignore-1 file)
	for _, issue := range issues {
		if issue.ID == "core-001" || issue.ID == "core-007" {
			t.Errorf("found ignored rule %s in output", issue.ID)
		}
	}
}

func TestExternalRulesFile(t *testing.T) {
	stdout, stderr, exitCode := runCLI("-r", "../../testdata/custom-rules.yaml", "-R", "none", "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("expected valid JSON: %v\nGot: %s", err, stdout)
	}

	// The custom rule detects EXPOSE directive
	foundCustom := false
	for _, issue := range issues {
		if issue.ID == "custom-001" {
			foundCustom = true
			break
		}
	}
	if !foundCustom {
		t.Error("expected to find custom-001 from external rules file")
	}
}

func TestCombinedFlags(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "combined-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Combine -q, -E, and -o
	_, stderr, exitCode := runCLI("-q", "-E", "-o", tmpFile.Name(), "../../testdata/Dockerfile-example")
	if exitCode != 1 {
		t.Errorf("expected exit code 1 with -E and issues, got %d, stderr: %s", exitCode, stderr)
	}

	// Output file should exist and have content
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	var issues []rules.Issue
	if err := json.Unmarshal(data, &issues); err != nil {
		t.Fatalf("expected valid JSON in output file: %v", err)
	}

	if len(issues) == 0 {
		t.Error("expected issues in output file")
	}
}

func TestMissingDockerfile(t *testing.T) {
	_, stderr, exitCode := runCLI("/nonexistent/Dockerfile")
	if exitCode == 0 {
		t.Error("expected non-zero exit code for missing file")
	}
	if !strings.Contains(stderr, "reading Dockerfile") && !strings.Contains(stderr, "no such file") {
		t.Errorf("expected error message about missing file, got: %s", stderr)
	}
}

func TestInvalidRuleSet(t *testing.T) {
	_, stderr, exitCode := runCLI("-R", "invalid", "../../testdata/Dockerfile-example")
	if exitCode == 0 {
		t.Error("expected non-zero exit code for invalid rule set")
	}
	if !strings.Contains(stderr, "unknown rule category") && !strings.Contains(stderr, "unknown rule selection") {
		t.Errorf("expected error about unknown rule selection or category, got: %s", stderr)
	}
}

func TestNoInputProvided(t *testing.T) {
	// Run without file and without stdin input
	cmd := exec.Command(binaryPath)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()

	if err == nil {
		t.Error("expected error when no input provided")
	}

	stderr := errBuf.String()
	if !strings.Contains(stderr, "Dockerfile is needed") {
		t.Errorf("expected 'Dockerfile is needed' error, got: %s", stderr)
	}
}

func TestMultipleIgnoreFlags(t *testing.T) {
	stdout, stderr, exitCode := runCLI("-i", "core-002", "-i", "core-003", "../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}

	// Should not contain core-002 or core-003
	for _, issue := range issues {
		if issue.ID == "core-002" || issue.ID == "core-003" {
			t.Errorf("found ignored rule %s in output", issue.ID)
		}
	}
}

func TestJSONOutputFormat(t *testing.T) {
	stdout, _, _ := runCLI("../../testdata/Dockerfile-example")

	var issues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &issues); err != nil {
		t.Fatalf("output should be valid JSON: %v\nGot: %s", err, stdout)
	}

	// Verify JSON structure
	for _, issue := range issues {
		if issue.ID == "" {
			t.Error("issue ID should not be empty")
		}
		if issue.Description == "" {
			t.Error("issue description should not be empty")
		}
		if issue.Severity == "" {
			t.Error("issue severity should not be empty")
		}
	}
}

func TestGoldenJSONOutput(t *testing.T) {
	stdout, stderr, exitCode := runCLI("../../testdata/Dockerfile-example")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", exitCode, stderr)
	}

	// Read expected output
	expected, err := os.ReadFile("../../testdata/golden/json_output.json")
	if err != nil {
		t.Fatalf("reading golden file: %v", err)
	}

	// Parse both to compare
	var actualIssues, expectedIssues []rules.Issue
	if err := json.Unmarshal([]byte(stdout), &actualIssues); err != nil {
		t.Fatalf("parsing actual output: %v", err)
	}
	if err := json.Unmarshal(expected, &expectedIssues); err != nil {
		t.Fatalf("parsing expected output: %v", err)
	}

	if len(actualIssues) != len(expectedIssues) {
		t.Errorf("issue count mismatch: got %d, want %d", len(actualIssues), len(expectedIssues))
	}

	// Build map of expected issues by ID
	expectedMap := make(map[string]rules.Issue)
	for _, issue := range expectedIssues {
		expectedMap[issue.ID] = issue
	}

	// Verify each actual issue matches expected
	for _, actual := range actualIssues {
		expected, ok := expectedMap[actual.ID]
		if !ok {
			t.Errorf("unexpected issue ID: %s", actual.ID)
			continue
		}
		if actual.Description != expected.Description {
			t.Errorf("issue %s description mismatch: got %q, want %q", actual.ID, actual.Description, expected.Description)
		}
		if actual.Severity != expected.Severity {
			t.Errorf("issue %s severity mismatch: got %q, want %q", actual.ID, actual.Severity, expected.Severity)
		}
	}
}
