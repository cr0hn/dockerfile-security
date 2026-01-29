package rules

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cr0hn/dockerfile-sec/internal/rules/embedded"
	"gopkg.in/yaml.v3"
)

// Rule represents a single security rule loaded from YAML.
type Rule struct {
	ID          string `yaml:"id" json:"id"`
	Description string `yaml:"description" json:"description"`
	Regex       string `yaml:"regex" json:"-"`
	Reference   string `yaml:"reference" json:"reference"`
	Severity    string `yaml:"severity" json:"severity"`
}

// Issue represents a matched rule (without the regex field).
type Issue struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Reference   string `json:"reference"`
	Severity    string `json:"severity"`
}

// IssueFromRule creates an Issue from a Rule, omitting the regex.
func IssueFromRule(r Rule) Issue {
	return Issue{
		ID:          r.ID,
		Description: r.Description,
		Reference:   r.Reference,
		Severity:    r.Severity,
	}
}

// LoadInternal loads built-in rules based on the selection flag.
// Valid selections: "all" (default), "core", "credentials", "security", "packages", "configuration", "none", or comma-separated combinations.
func LoadInternal(selection string) ([]Rule, error) {
	selection = strings.ToLower(selection)

	// Handle special cases
	if selection == "none" {
		return nil, nil
	}

	if selection == "all" || selection == "" {
		return loadAllCategories()
	}

	// Handle comma-separated categories
	if strings.Contains(selection, ",") {
		categories := strings.Split(selection, ",")
		var allRules []Rule
		for _, cat := range categories {
			cat = strings.TrimSpace(cat)
			rules, err := loadCategory(cat)
			if err != nil {
				return nil, err
			}
			allRules = append(allRules, rules...)
		}
		return allRules, nil
	}

	// Single category
	return loadCategory(selection)
}

// loadCategory loads a single category of rules.
func loadCategory(category string) ([]Rule, error) {
	switch category {
	case "core":
		return parseYAML(embedded.CoreYAML)
	case "credentials":
		return parseYAML(embedded.CredentialsYAML)
	case "security":
		return parseYAML(embedded.SecurityYAML)
	case "packages":
		return parseYAML(embedded.PackagesYAML)
	case "configuration":
		return parseYAML(embedded.ConfigurationYAML)
	default:
		return nil, fmt.Errorf("unknown rule category: %s", category)
	}
}

// loadAllCategories loads all built-in rule categories.
func loadAllCategories() ([]Rule, error) {
	var allRules []Rule

	categories := [][]byte{
		embedded.CoreYAML,
		embedded.CredentialsYAML,
		embedded.SecurityYAML,
		embedded.PackagesYAML,
		embedded.ConfigurationYAML,
	}

	for _, data := range categories {
		rules, err := parseYAML(data)
		if err != nil {
			return nil, err
		}
		allRules = append(allRules, rules...)
	}

	return allRules, nil
}

// LoadExternal loads rules from a file path or URL.
func LoadExternal(source string) ([]Rule, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return LoadFromURL(source)
	}
	return LoadFromFile(source)
}

// LoadFromFile loads rules from a local YAML file.
func LoadFromFile(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading rules file %s: %w", path, err)
	}
	return parseYAML(data)
}

// LoadFromURL loads rules from a remote YAML URL.
func LoadFromURL(url string) ([]Rule, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching rules from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching rules from %s: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading rules from %s: %w", url, err)
	}
	return parseYAML(data)
}

func parseYAML(data []byte) ([]Rule, error) {
	var rules []Rule
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("parsing rules YAML: %w", err)
	}
	return rules, nil
}
