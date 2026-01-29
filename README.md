<p align="center">
  <img src="https://raw.githubusercontent.com/cr0hn/dockerfile-security/master/docs/logo.png" alt="dockerfile-sec logo" width="200"/>
</p>

<h1 align="center">dockerfile-sec</h1>

<p align="center">
  <strong>A fast, rule-based security scanner for Dockerfiles</strong>
</p>

<p align="center">
  Detect misconfigurations, exposed credentials, and security anti-patterns before they reach production.
</p>

<p align="center">
  <a href="https://github.com/cr0hn/dockerfile-security/actions/workflows/ci.yml"><img src="https://github.com/cr0hn/dockerfile-security/actions/workflows/ci.yml/badge.svg" alt="CI/CD"></a>
  <a href="https://goreportcard.com/report/github.com/cr0hn/dockerfile-security"><img src="https://goreportcard.com/badge/github.com/cr0hn/dockerfile-security" alt="Go Report Card"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-BSD_3--Clause-blue.svg" alt="License"></a>
  <a href="https://github.com/cr0hn/dockerfile-security/pkgs/container/dockerfile-security"><img src="https://img.shields.io/badge/Docker-ghcr.io-2496ED?style=flat&logo=docker" alt="Docker"></a>
  <a href="https://github.com/cr0hn/dockerfile-security/releases"><img src="https://img.shields.io/github/v/release/cr0hn/dockerfile-security?sort=semver" alt="Latest Release"></a>
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#built-in-rules">Rules</a> •
  <a href="#creating-custom-rules">Custom Rules</a> •
  <a href="#contributing">Contributing</a>
</p>

---

## Table of Contents

- [Why dockerfile-sec?](#why-dockerfile-sec)
- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
  - [Binary](#binary)
  - [Docker](#docker)
  - [From Source](#from-source)
- [Usage](#usage)
  - [Basic Analysis](#basic-analysis)
  - [Using Docker](#using-docker)
  - [Pipeline Integration](#pipeline-integration)
  - [CI/CD Integration](#cicd-integration)
- [Configuration](#configuration)
  - [Rule Sets](#rule-sets)
  - [Output Formats](#output-formats)
  - [Ignoring Rules](#ignoring-rules)
  - [External Rules](#external-rules)
- [Built-in Rules](#built-in-rules)
  - [Core Rules](#core-rules)
  - [Credential Rules](#credential-rules)
- [Creating Custom Rules](#creating-custom-rules)
  - [Rule Format](#rule-format)
  - [Regex Tips](#regex-tips)
  - [Examples](#rule-examples)
- [CLI Reference](#cli-reference)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)
- [Acknowledgments](#acknowledgments)

---

## Why dockerfile-sec?

Dockerfiles can contain security issues that are easy to miss during code reviews:

- **Hardcoded credentials** - Passwords, API keys, and tokens accidentally committed
- **Running as root** - Missing `USER` directive leads to container privilege escalation
- **Insecure base images** - Using `latest` tag or images without SHA256 verification
- **Exposed secrets in build args** - Sensitive data passed via `ARG` instead of secrets
- **Recursive copies** - `COPY . .` accidentally including `.env` files and credentials

**dockerfile-sec** catches these issues automatically, integrating seamlessly into your development workflow and CI/CD pipelines.

---

## Features

| Feature | Description |
|---------|-------------|
| **35 Built-in Rules** | Comprehensive coverage of security best practices and credential detection |
| **Blazing Fast** | Written in Go for maximum performance on large codebases |
| **Flexible Output** | ASCII tables for humans, JSON for machines and automation |
| **CI/CD Ready** | Exit codes and quiet mode for seamless pipeline integration |
| **Extensible** | Load custom rules from local files or remote URLs |
| **Zero Dependencies** | Single static binary, no runtime required |
| **Docker Support** | Available as a minimal container image |
| **Cross-Platform** | Linux, macOS, and Windows support (amd64/arm64) |

---

## Quick Start

```bash
# Download the binary (Linux/amd64)
curl -L https://github.com/cr0hn/dockerfile-security/releases/latest/download/dockerfile-sec-linux-amd64 -o dockerfile-sec
chmod +x dockerfile-sec

# Scan a Dockerfile
./dockerfile-sec Dockerfile

# Scan with exit code for CI/CD (exits 1 if issues found)
./dockerfile-sec -E Dockerfile
```

**Example output:**

```
+----------+-------------------------------------------+----------+
| Rule Id  | Description                               | Severity |
+----------+-------------------------------------------+----------+
| core-002 | Posible text plain password in dockerfile | High     |
| core-003 | Recursive copy found                      | Medium   |
| core-005 | Use image tag instead of SHA256 hash      | Medium   |
| cred-001 | Generic credential                        | Medium   |
+----------+-------------------------------------------+----------+
```

---

## Installation

### Binary

Download the latest release for your platform from the [releases page](https://github.com/cr0hn/dockerfile-security/releases).

| Platform | Architecture | Download |
|----------|--------------|----------|
| Linux    | amd64        | [dockerfile-sec-linux-amd64](https://github.com/cr0hn/dockerfile-security/releases/latest/download/dockerfile-sec-linux-amd64) |
| Linux    | arm64        | [dockerfile-sec-linux-arm64](https://github.com/cr0hn/dockerfile-security/releases/latest/download/dockerfile-sec-linux-arm64) |
| macOS    | amd64        | [dockerfile-sec-darwin-amd64](https://github.com/cr0hn/dockerfile-security/releases/latest/download/dockerfile-sec-darwin-amd64) |
| macOS    | arm64 (M1/M2)| [dockerfile-sec-darwin-arm64](https://github.com/cr0hn/dockerfile-security/releases/latest/download/dockerfile-sec-darwin-arm64) |
| Windows  | amd64        | [dockerfile-sec-windows-amd64.exe](https://github.com/cr0hn/dockerfile-security/releases/latest/download/dockerfile-sec-windows-amd64.exe) |

### Docker

```bash
# Using GitHub Container Registry
docker pull ghcr.io/cr0hn/dockerfile-security:latest

# Scan a Dockerfile via stdin
cat Dockerfile | docker run --rm -i ghcr.io/cr0hn/dockerfile-security

# Scan a local file (mount as volume)
docker run --rm -v $(pwd):/app ghcr.io/cr0hn/dockerfile-security /app/Dockerfile
```

### From Source

Requires Go 1.22 or later.

```bash
# Install directly
go install github.com/cr0hn/dockerfile-security/cmd/dockerfile-sec@latest

# Or build manually
git clone https://github.com/cr0hn/dockerfile-security.git
cd dockerfile-security
make build
```

---

## Usage

### Basic Analysis

```bash
# Scan a Dockerfile (outputs ASCII table in terminal, JSON when piped)
dockerfile-sec Dockerfile

# Read from stdin
cat Dockerfile | dockerfile-sec

# Quiet mode with exit code for scripts
if dockerfile-sec -E -q Dockerfile; then
  echo "No security issues found"
else
  echo "Security issues detected!"
  exit 1
fi
```

### Using Docker

The Docker image is available at `ghcr.io/cr0hn/dockerfile-security`.

**Basic usage:**

```bash
# Scan a Dockerfile via stdin
cat Dockerfile | docker run --rm -i ghcr.io/cr0hn/dockerfile-security

# Scan a local file (mount current directory)
docker run --rm -v $(pwd):/workspace ghcr.io/cr0hn/dockerfile-security /workspace/Dockerfile

# Scan with exit code for CI/CD
docker run --rm -v $(pwd):/workspace ghcr.io/cr0hn/dockerfile-security -E /workspace/Dockerfile
```

**Advanced usage:**

```bash
# Use a specific version
docker run --rm -i ghcr.io/cr0hn/dockerfile-security:v1.0.0

# With custom rules (mount rules file)
docker run --rm -i \
  -v $(pwd)/Dockerfile:/Dockerfile \
  -v $(pwd)/custom-rules.yaml:/rules.yaml \
  ghcr.io/cr0hn/dockerfile-security -r /rules.yaml /Dockerfile

# Output JSON to a file
docker run --rm -v $(pwd):/workspace ghcr.io/cr0hn/dockerfile-security \
  -o /workspace/results.json /workspace/Dockerfile

# Ignore specific rules
docker run --rm -v $(pwd):/workspace ghcr.io/cr0hn/dockerfile-security \
  -i core-001 -i core-004 /workspace/Dockerfile

# Only credential rules
docker run --rm -v $(pwd):/workspace ghcr.io/cr0hn/dockerfile-security \
  -R credentials /workspace/Dockerfile
```

**Docker Compose integration:**

```yaml
# docker-compose.yml
services:
  dockerfile-sec:
    image: ghcr.io/cr0hn/dockerfile-security:latest
    volumes:
      - ./Dockerfile:/Dockerfile:ro
    command: ["-E", "/Dockerfile"]
```

**Shell alias for convenience:**

```bash
# Add to ~/.bashrc or ~/.zshrc
alias dockerfile-sec='docker run --rm -i -v $(pwd):/workspace ghcr.io/cr0hn/dockerfile-security'

# Usage
dockerfile-sec /workspace/Dockerfile
dockerfile-sec -E /workspace/Dockerfile
```

### Pipeline Integration

dockerfile-sec works seamlessly in UNIX pipelines:

```bash
# Chain with jq for JSON processing
cat Dockerfile | dockerfile-sec | jq '.[] | select(.severity == "High")'

# Process multiple Dockerfiles
find . -name "Dockerfile*" -exec dockerfile-sec -E {} \;
```

### CI/CD Integration

#### GitHub Actions

```yaml
name: Security Scan

on: [push, pull_request]

jobs:
  dockerfile-security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Scan Dockerfile
        run: |
          curl -sL https://github.com/cr0hn/dockerfile-security/releases/latest/download/dockerfile-sec-linux-amd64 -o dockerfile-sec
          chmod +x dockerfile-sec
          ./dockerfile-sec -E Dockerfile
```

#### GitLab CI

```yaml
dockerfile-security:
  image: ghcr.io/cr0hn/dockerfile-security:latest
  script:
    - dockerfile-sec -E Dockerfile
  rules:
    - changes:
        - Dockerfile
        - "*.dockerfile"
```

#### Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('Dockerfile Security') {
            steps {
                sh '''
                    curl -sL https://github.com/cr0hn/dockerfile-security/releases/latest/download/dockerfile-sec-linux-amd64 -o dockerfile-sec
                    chmod +x dockerfile-sec
                    ./dockerfile-sec -E Dockerfile
                '''
            }
        }
    }
}
```

---

## Configuration

### Rule Sets

Control which built-in rules are loaded:

```bash
# All rules (default)
dockerfile-sec Dockerfile

# Core rules only (best practices)
dockerfile-sec -R core Dockerfile

# Credential detection rules only
dockerfile-sec -R credentials Dockerfile

# Security rules only
dockerfile-sec -R security Dockerfile

# Package management rules only
dockerfile-sec -R packages Dockerfile

# Configuration rules only
dockerfile-sec -R configuration Dockerfile

# Combine multiple categories (comma-separated)
dockerfile-sec -R core,security Dockerfile
dockerfile-sec -R credentials,security,packages Dockerfile

# Disable built-in rules (use with -r for custom rules only)
dockerfile-sec -R none -r my-rules.yaml Dockerfile
```

### Output Formats

```bash
# ASCII table (default in terminal)
dockerfile-sec Dockerfile

# JSON output (automatic when piped, or explicit with -o)
dockerfile-sec Dockerfile | cat
dockerfile-sec -o results.json Dockerfile

# Quiet mode (no output, useful with -E for CI/CD)
dockerfile-sec -q -E Dockerfile
```

**JSON Output Format:**

```json
[
  {
    "id": "core-002",
    "description": "Posible text plain password in dockerfile",
    "reference": "https://snyk.io/blog/10-docker-image-security-best-practices/",
    "severity": "High"
  }
]
```

### Ignoring Rules

**By rule ID (CLI):**

```bash
# Ignore single rule
dockerfile-sec -i core-001 Dockerfile

# Ignore multiple rules
dockerfile-sec -i core-001 -i core-007 Dockerfile
```

**By ignore file:**

Create a file with rule IDs to ignore (one per line):

```text
# .dockerfile-sec-ignore
# Ignore USER requirement for this project
core-001

# We use ADD intentionally for tar extraction
core-004
```

```bash
dockerfile-sec -F .dockerfile-sec-ignore Dockerfile
```

### External Rules

Load custom rules from files or URLs:

```bash
# From local file
dockerfile-sec -r my-rules.yaml Dockerfile

# From URL
dockerfile-sec -r https://example.com/rules.yaml Dockerfile

# Combine with built-in rules
dockerfile-sec -r my-rules.yaml Dockerfile

# Use only external rules
dockerfile-sec -R none -r my-rules.yaml Dockerfile
```

---

## Built-in Rules

dockerfile-sec includes **35 built-in rules** across 5 categories:

### Core Rules (10 rules)

Best practices and security guidelines for Dockerfiles.

| ID | Description | Severity |
|----|-------------|----------|
| `core-001` | Missing USER sentence (running as root) | High |
| `core-002` | Possible plaintext password in Dockerfile | High |
| `core-003` | Recursive copy found (`COPY . .`) | Medium |
| `core-004` | Use of ADD instead of COPY | Low |
| `core-005` | Use image tag instead of SHA256 hash | Medium |
| `core-006` | Use of `latest` tag in FROM | Medium |
| `core-007` | Use of deprecated MAINTAINER | Low |
| `core-008` | Use of `--insecurity=insecure` in RUN | High |
| `core-009` | Secrets passed via ARG instead of ENV | High |
| `core-010` | HEALTHCHECK contains sensitive information | High |

### Credential Rules (11 rules)

Detection of exposed secrets and credentials.

| ID | Description | Severity |
|----|-------------|----------|
| `cred-001` | Generic credential patterns | Medium |
| `cred-002` | AWS Access Key ID | High |
| `cred-003` | AWS MWS Key | High |
| `cred-004` | EC Private Key | High |
| `cred-005` | Google API Key | High |
| `cred-006` | Slack Webhook URL | High |
| `cred-007` | GitHub Personal Access Token | Critical |
| `cred-008` | RSA Private Key | Critical |
| `cred-009` | OpenAI API Key | High |
| `cred-010` | Stripe API Key | Critical |
| `cred-011` | Docker Registry Authentication Token | High |

### Security Rules (7 rules)

Advanced security checks for container security.

| ID | Description | Severity |
|----|-------------|----------|
| `sec-001` | Docker socket mounted in container (container escape) | Critical |
| `sec-002` | WORKDIR pointing to system directories or root | Medium |
| `sec-003` | chmod 777 or dangerous permissions | High |
| `sec-004` | Use of sudo in RUN commands | Medium |
| `sec-005` | BuildKit secret mount without cleanup | High |
| `sec-006` | Setting SUID/SGID bits on binaries | High |
| `sec-007` | ENV directive with embedded credentials | High |

### Package Rules (4 rules)

Package manager best practices and security.

| ID | Description | Severity |
|----|-------------|----------|
| `pkg-001` | apt-get without cleanup | Medium |
| `pkg-002` | pip install without --no-cache-dir | Low |
| `pkg-003` | npm install without cache cleanup | Low |
| `pkg-004` | Piping curl/wget to bash | High |

### Configuration Rules (3 rules)

Container runtime configuration issues.

| ID | Description | Severity |
|----|-------------|----------|
| `cfg-001` | Using --privileged flag | Critical |
| `cfg-002` | Exposing dangerous ports (22, 23, 3389, etc.) | Medium |
| `cfg-003` | Non-standard STOPSIGNAL defined | Low |

---

## Creating Custom Rules

### Rule Format

Rules are defined in YAML format. Each rule requires:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier (e.g., `custom-001`) |
| `description` | string | Yes | Human-readable description |
| `regex` | string | Yes | Regular expression pattern to match |
| `reference` | string | Yes | URL with more information |
| `severity` | string | Yes | `Low`, `Medium`, or `High` |

### Rule Examples

**custom-rules.yaml:**

```yaml
# Detect hardcoded IP addresses
- id: custom-001
  description: Hardcoded IP address found
  regex: '(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})'
  reference: https://example.com/security-guidelines
  severity: Medium

# Detect curl without certificate verification
- id: custom-002
  description: Use of curl without certificate verification
  regex: '(curl.*-k|curl.*--insecure)'
  reference: https://example.com/security-guidelines
  severity: High

# Detect wget without certificate verification
- id: custom-003
  description: Use of wget without certificate verification
  regex: '(wget.*--no-check-certificate)'
  reference: https://example.com/security-guidelines
  severity: High

# Detect EXPOSE directive
- id: custom-004
  description: EXPOSE directive detected - verify port is necessary
  regex: '^(EXPOSE[\s]+[\d]+)'
  reference: https://docs.docker.com/reference/dockerfile/#expose
  severity: Low
```

**Using custom rules:**

```bash
# Add to built-in rules
dockerfile-sec -r custom-rules.yaml Dockerfile

# Replace built-in rules entirely
dockerfile-sec -R none -r custom-rules.yaml Dockerfile

# Multiple rule files
dockerfile-sec -r rules1.yaml -r rules2.yaml Dockerfile
```

### Regex Tips

| Tip | Example |
|-----|---------|
| Match at line start | `^FROM` |
| Case-insensitive | `(?i)password` |
| Match any character | `.` |
| Match literal dot | `\.` |
| Match word boundary | `\b` |
| Non-greedy match | `.*?` |

Test your regex at [regex101.com](https://regex101.com/) (select Go flavor).

---

## CLI Reference

```
Usage: dockerfile-sec [OPTIONS] [DOCKERFILE]

Analyze a Dockerfile for security issues.

Arguments:
  DOCKERFILE    Path to Dockerfile (reads from stdin if not provided)

Options:
  -E            Exit with code 1 if issues are found (for CI/CD)
  -F file       Ignore file containing rule IDs to skip (repeatable)
  -R selection  Built-in rules: all, core, credentials, security, packages, configuration, none (comma-separated, default: all)
  -i id         Ignore specific rule ID (repeatable)
  -o file       Write JSON output to file
  -q            Quiet mode (suppress stdout output)
  -r file       External rules file or URL (repeatable)
  -h, --help    Show help message
  -v, --version Show version information

Exit Codes:
  0             Success (no issues found, or -E not specified)
  1             Issues found (only when -E is specified)
  2             Error (invalid arguments, file not found, etc.)
```

---

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Setting up the development environment
- Code style guidelines
- How to submit pull requests
- Adding new security rules

Before contributing, please read our [Code of Conduct](CODE_OF_CONDUCT.md).

### Quick Start for Contributors

```bash
# Clone the repository
git clone https://github.com/cr0hn/dockerfile-security.git
cd dockerfile-security

# Run tests
make test

# Run linter
make lint

# Build
make build
```

---

## Security

For information about reporting security vulnerabilities, please see our [Security Policy](SECURITY.md).

**Please do not report security vulnerabilities through public GitHub issues.**

---

## License

This project is licensed under the **BSD 3-Clause License** - see the [LICENSE](LICENSE) file for details.

```
BSD 3-Clause License

Copyright (c) 2020-2025, Daniel Garcia (cr0hn)
All rights reserved.
```

---

## Acknowledgments

- [Snyk](https://snyk.io/blog/10-docker-image-security-best-practices/) - Docker security best practices
- [gitleaks](https://github.com/zricethezav/gitleaks) - Credential detection patterns
- [Docker Documentation](https://docs.docker.com/reference/dockerfile/) - Dockerfile reference

---

## References

- [Snyk: 10 Docker Image Security Best Practices](https://snyk.io/blog/10-docker-image-security-best-practices/)
- [Dockerfile Security Tuneup](https://medium.com/microscaling-systems/dockerfile-security-tuneup-166f1cdafea1)
- [Container Deployments: A Lesson in Deterministic Ops](https://medium.com/@tariq.m.islam/container-deployments-a-lesson-in-deterministic-ops-a4a467b14a03)
- [Spacelift: Docker Security Best Practices](https://spacelift.io/blog/docker-security)
- [Docker Documentation: Dockerfile Reference](https://docs.docker.com/reference/dockerfile/)

---

## Related Projects

You might also be interested in **[dockerscan](https://github.com/cr0hn/dockerscan)** - A Docker image security analyzer that complements dockerfile-sec by scanning built Docker images for vulnerabilities, malware, and security issues.

---

<p align="center">
  <sub>Made with ❤️ by <a href="https://github.com/cr0hn">Daniel Garcia (cr0hn)</a></sub>
</p>

<p align="center">
  <a href="https://github.com/cr0hn/dockerfile-security/stargazers">⭐ Star us on GitHub</a>
</p>
