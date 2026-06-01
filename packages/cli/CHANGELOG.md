# Changelog

All notable changes to AutoDev will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added

- Initial monorepo structure with Go workspace and PNPM workspaces
- CLI scaffold with Cobra commands: scan, setup, github, doctor, report, install, update, clean, skills, export
- Repository scanner for 30+ languages and frameworks
- GitHub user repository scanner
- Skills.sh integration for learning roadmaps
- Neo-Brutalist website with Next.js 15
- Docker support
- GoReleaser configuration for multi-platform binaries
- GitHub Actions CI/CD pipeline
- curl-pipe installer script

---

## [0.1.0] - 2026-05-29

### Added

- Initial release 🎉
- Core CLI with `autodev scan` and `autodev setup`
- Support for Node.js, Python, Go, Rust, Java, Docker detection
- HTML, JSON, and Markdown report generation
- Multi-platform binaries (Linux x86_64/arm64, macOS, Windows)

[Unreleased]: https://github.com/HEETMEHTA18/autodev/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/HEETMEHTA18/autodev/releases/tag/v0.1.0
