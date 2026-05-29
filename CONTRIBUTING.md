# Contributing to AutoDev

Thank you for your interest in contributing to AutoDev! 🎉

We welcome contributions from developers of all experience levels.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Reporting Issues](#reporting-issues)

---

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

---

## Getting Started

### Prerequisites

- **Go** ≥ 1.22
- **Node.js** ≥ 20
- **pnpm** ≥ 9
- **Git**

### Fork and Clone

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/YOUR_USERNAME/autodev.git
cd autodev
```

---

## Development Setup

### Install all dependencies

```bash
# Install Node dependencies
pnpm install

# Download Go dependencies
cd packages/cli && go mod download
cd ../core && go mod download
cd ../scanner && go mod download
cd ../installer && go mod download
cd ../skills && go mod download
cd ../github && go mod download
cd ../..
```

### Build the CLI

```bash
# Fast build for development
pnpm cli:build

# Or directly with Go
cd packages/cli
go build -o bin/autodev ./main.go
./bin/autodev --help
```

### Run the website locally

```bash
pnpm website
# Opens http://localhost:3000
```

---

## Project Structure

```
autodev/
├── apps/
│   └── website/          # Next.js 15 marketing site
├── packages/
│   ├── cli/              # Main CLI entry point (Cobra + BubbleTea)
│   ├── core/             # Shared utilities (OS, config, logging)
│   ├── scanner/          # Repository detection engine
│   ├── installer/        # Runtime installer logic
│   ├── skills/           # skills.sh learning roadmap integration
│   └── github/           # GitHub API client
├── scripts/
│   ├── install.sh        # Curl-pipe installer
│   └── build.sh
├── .github/workflows/    # GitHub Actions
└── go.work               # Go workspace
```

---

## Making Changes

### Branching strategy

```bash
# For features
git checkout -b feat/your-feature-name

# For bug fixes
git checkout -b fix/bug-description

# For documentation
git checkout -b docs/topic
```

### Commit conventions

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add kotlin detection to scanner
fix: correct go binary path on windows
docs: update installation instructions
test: add scanner unit tests for rust projects
chore: update dependencies
```

---

## Submitting a Pull Request

1. Ensure all tests pass: `pnpm test` and `pnpm cli:test`
2. Update documentation if needed
3. Add/update tests for your changes
4. Open a PR with a clear title and description
5. Link any related issues with `Fixes #123`

---

## Coding Standards

### Go

- Follow standard Go formatting (`gofmt`)
- Use `golangci-lint` for linting
- Write tests for all exported functions
- Use meaningful variable names
- Add comments for exported types and functions

### TypeScript / Next.js

- Use TypeScript strictly (no `any`)
- Follow ESLint rules
- Use Prettier for formatting

---

## Testing

### CLI tests

```bash
cd packages/cli
go test ./... -v -race -cover
```

### Website tests

```bash
pnpm --filter=website test
```

### Coverage requirements

We require **≥ 90% coverage** for Go packages.

---

## Reporting Issues

- **Bugs**: Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md)
- **Features**: Use the [feature request template](.github/ISSUE_TEMPLATE/feature_request.md)
- **Security**: See [SECURITY.md](SECURITY.md) — **do NOT open a public issue**

---

## Need Help?

- Open a [Discussion](https://github.com/autodev-sh/autodev/discussions)
- Join our community on Discord: https://discord.gg/autodev

Thank you for contributing! 🚀
