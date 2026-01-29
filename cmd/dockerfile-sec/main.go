package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cr0hn/dockerfile-sec/internal/analyzer"
	"github.com/cr0hn/dockerfile-sec/internal/ignore"
	"github.com/cr0hn/dockerfile-sec/internal/output"
	"github.com/cr0hn/dockerfile-sec/internal/rules"
)

// stringSliceFlag implements flag.Value for repeatable string flags.
type stringSliceFlag []string

func (s *stringSliceFlag) String() string { return strings.Join(*s, ",") }
func (s *stringSliceFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "[!]  %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		ignoreFiles   stringSliceFlag
		ignoreRules   stringSliceFlag
		rulesFiles    stringSliceFlag
		internalRules string
		outputFile    string
		quiet         bool
		codeExit      bool
	)

	flag.Var(&ignoreFiles, "F", "ignore file (repeatable)")
	flag.Var(&ignoreRules, "i", "ignore rule IDs, comma-separated (repeatable)")
	flag.Var(&rulesFiles, "r", "external rules file or URL (repeatable)")
	flag.StringVar(&internalRules, "R", "all", "built-in rules: core, credentials, security, packages, configuration, all, none (comma-separated)")
	flag.StringVar(&outputFile, "o", "", "output file path (JSON)")
	flag.BoolVar(&quiet, "q", false, "quiet mode")
	flag.BoolVar(&codeExit, "E", false, "exit code 1 if issues found")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dockerfile-sec [OPTIONS] [DOCKERFILE]\n\nAnalyze a Dockerfile for security issues.\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Load Dockerfile content
	var content string
	args := flag.Args()

	if len(args) > 0 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading Dockerfile: %w", err)
		}
		content = string(data)
	} else {
		// Try reading from stdin
		info, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("checking stdin: %w", err)
		}
		if info.Mode()&os.ModeCharDevice != 0 {
			return fmt.Errorf("Dockerfile is needed")
		}
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
		content = string(data)
	}

	if content == "" {
		return fmt.Errorf("Dockerfile is needed")
	}

	// Load rules
	allRules, err := rules.LoadInternal(internalRules)
	if err != nil {
		return err
	}

	for _, rf := range rulesFiles {
		ext, err := rules.LoadExternal(rf)
		if err != nil {
			return err
		}
		allRules = append(allRules, ext...)
	}

	// Load ignores
	ignored, err := ignore.Load(ignoreRules, ignoreFiles)
	if err != nil {
		return err
	}

	// Analyze
	issues := analyzer.Analyze(content, allRules, ignored)

	// Output
	if err := output.Render(issues, quiet, outputFile); err != nil {
		return err
	}

	// Exit code
	if codeExit && len(issues) > 0 {
		os.Exit(1)
	}

	return nil
}
