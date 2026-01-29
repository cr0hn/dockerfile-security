# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2025-01-29

### Added

- **19 New Security Rules** - Expanded from 16 to 35 built-in rules
- **New Rule Categories**:
  - `security` (7 rules) - Advanced security checks including Docker socket mounts, dangerous permissions, SUID/SGID bits
  - `packages` (4 rules) - Package manager best practices for apt, pip, npm, and curl|bash detection
  - `configuration` (3 rules) - Runtime configuration issues like --privileged flag and dangerous ports
- **5 New Credential Detection Rules** - GitHub PAT, RSA Private Key, OpenAI API Key, Stripe API Key, Docker Registry Auth
- **Comma-separated rule categories** - Load multiple categories at once: `-R core,security,packages`
- Enhanced core rules with improved regex patterns and severity levels

### Changed

- `core-001` severity upgraded from Medium to High (missing USER directive)
- `core-002` improved regex to detect more password patterns
- `core-005` simplified regex for SHA256 hash detection
- `core-006` simplified regex for latest tag detection
- `core-009` expanded keywords for better secret detection
- Rule count: 16 → 35 (10 core + 11 credentials + 7 security + 4 packages + 3 configuration)

### Fixed

- Fixed `-R` flag help text to include new categories
- Updated test suite to validate all 35 rules
- Updated golden files to reflect new rule detections

## [0.2.0] - 2025-01-29

### Changed

- **Complete rewrite from Python to Go 1.22** - Full language migration while maintaining 100% CLI compatibility
- Modular architecture with clean separation of concerns (`cmd/`, `internal/analyzer/`, `internal/rules/`, `internal/output/`, `internal/ignore/`)
- Upgraded regex engine from Python `re` to `regexp2` with 5-second timeout protection
- Improved error handling and reporting throughout the codebase
- Static binary distribution (zero runtime dependencies, no CGO)

### Added

- **Comprehensive test suite** with 1,071 lines of tests covering:
  - 20 end-to-end CLI tests
  - 6 analyzer tests
  - 15 rule loading tests
  - 9 output formatter tests
  - 5 ignore system tests
  - Test coverage: 80%+ across all packages
- **CI/CD pipeline** with GitHub Actions:
  - Automated testing with race detection
  - Linting with `go vet` and `gofmt`
  - Multi-platform Docker image builds
  - Cross-platform binary releases (Linux, macOS, Windows for amd64/arm64)
- **Embedded rules** using `go:embed` - rules now bundled in the binary
- **Makefile** with build, test, coverage, and formatting targets
- **Build-time version injection** with commit hash tracking
- **CONTRIBUTING.md** with development setup and guidelines
- **Golden file testing** for output consistency verification
- Support for Apple Silicon (arm64 macOS binaries)

### Fixed

- Exit code behavior now properly returns 1 when `-E` flag is used and issues are found
- Better handling of edge cases (missing files, invalid rules, empty input)
- Regex validation with timeout protection against ReDoS attacks

### Removed

- Python 3.8+ runtime dependency
- PyPI package dependencies (`pyyaml`, `terminaltables`)
- Monolithic `dockerfile_sec/__main__.py` implementation

### Distribution

- Docker images now available at `ghcr.io/cr0hn/dockerfile-security`
- Pre-built binaries for all major platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64/M1/M2)
  - Windows (amd64)
- Single static binary with no external dependencies

### Performance

- Near-instantaneous startup (static binary vs. Python interpreter)
- Minimal memory footprint
- Faster regex matching with .NET-compatible engine

### Compatibility

- **Fully backward compatible** with v0.1.x CLI interface
- All flags preserved: `-F`, `-i`, `-r`, `-R`, `-o`, `-q`, `-E`
- Identical output formats (ASCII table, JSON)
- Same exit code behavior

## [0.1.0] - 2024

### Added

- Initial Python implementation
- 16 built-in security rules (10 core, 6 credential detection)
- YAML-based rule system
- ASCII table and JSON output formats
- Support for external rules from files and URLs
- Rule ignore system via CLI flags and files
- Exit code mode for CI/CD integration
- Stdin input support
- PyPI package distribution

---

**Note:** Version 0.2.0 represents a complete rewrite to Go while maintaining full compatibility with the 0.1.x Python implementation. All features have been preserved and enhanced with better performance, testing, and distribution.
