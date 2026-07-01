# golc — CLI Line Counter

A lightweight command-line wrapper around GoLC's analysis engine. Produces formatted terminal output similar to [scc](https://github.com/boyter/scc).

## Why this exists

Our organisation uses Sonar for code scanning and needs a practical way to estimate lines of code for onboarding projects before they are fully brought into the platform. In practice, we found GoLC to be the closest available approximation, with some remaining margin of error, for the LOC Sonar is likely to count and bill against.

The upstream GoLC project is geared around its existing web-based flow and does not currently provide a simple CLI entry point for the `golc [directory]` workflow we need. This wrapper was added to make local directory analysis usable from the command line.

All custom changes were intentionally added in a non-interfering way: new CLI files live under `cmd/golc-cli`, CI was added as a separate workflow, and the upstream project files remain untouched where possible. That keeps it straightforward to continue accepting upstream changes and PRs from the original Sonar repository so this fork stays aligned with Sonar's behaviour over time.

## Install

```bash
make install
```

This builds the binary and places it in `~/.local/bin/golc`. Ensure `~/.local/bin` is on your `PATH`.

To install elsewhere:

```bash
make install INSTALL_DIR=/usr/local/bin
```

## Usage

```bash
golc [directory]
```

If no directory is given, the current directory is scanned.

### Examples

```bash
# Scan current directory
golc

# Scan a specific project
golc ~/Projects/my-app

# Scan a subdirectory
golc ./src
```

### Sample output

```text
   Language  | Files | Lines | Blank lines | Comments | Code lines
-------------+-------+-------+-------------+----------+-------------
  Golang     |    60 | 22873 |        2859 |     1643 |      18371
  JavaScript |    17 | 10551 |         992 |      394 |       9165
  CSS        |    16 | 32550 |        7233 |      217 |      25100
  JSON       |     2 |   225 |           0 |        0 |        225
  Shell      |     5 |   156 |          31 |       21 |        104
  YAML       |     1 |    76 |          14 |        0 |         62
-------------+-------+-------+-------------+----------+-------------
    Total    |  101  | 66431 |    11129    |   2275   |   53027
```

## Build

```bash
# Build for current platform
make build

# Cross-compile for all platforms
make build-all

# Output binaries
ls bin/
# golc-linux-amd64  golc-darwin-arm64  golc-windows-amd64.exe  ...
```

## Test

```bash
# All tests (package + CLI + integration)
make test

# CLI tests only
make test-cli
```

## Supported languages

All languages supported by [SonarQube](https://www.sonarsource.com/knowledge/languages/) — including Go, Java, Python, TypeScript, C#, C/C++, Kotlin, Ruby, Rust, Swift, Terraform, and more. See the full list in the main project README.

## How it works

This CLI is a thin wrapper around the GoLC analysis engine (`pkg/goloc`). It:

1. Resolves the target directory
2. Walks the file tree, matching files by extension against SonarQube-supported languages
3. Counts physical lines, blank lines, comments, and code lines per file
4. Aggregates by language and prints a formatted table to stdout

No config file, no network access, no cloning — just local file analysis.

## Uninstall

```bash
make uninstall
```
