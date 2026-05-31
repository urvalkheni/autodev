# Changelog

All notable changes to AutoDev will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

---

## [0.2.0] - 2026-05-31

### Added
- **Ecosystem Expansion**: Added detection and automated setup for Docker and Kubernetes (Helm, kubectl).
- **Mobile/SDK Toolchains**: Native Flutter SDK installer, Android SDK, and emulator setup.
- **JVM & Build Systems**: OpenJDK installer (Java 17/21), Kotlin build tools, and Apache Maven.
- **AI Rule Files & 99.8% Context Saving**: Automatic generation of standard AI assistant rule files (`.autodev-skills.md`, `.cursorrules`, `.clinerules`, and `.github/copilot-instructions.md`) when running any AutoDev command. Reduces context payload by 99.8% (~350 tokens instead of 200,000+).
- **Interactive Visual Roadmap TUI**: A beautiful BubbleTea visual terminal interface for navigating skill roadmaps.
- **Model Context Protocol (MCP) Integration**: Native MCP server (`autodev mcp start` and auto-configuration via `autodev mcp setup`) exposing `autodev_scan`, `autodev_doctor`, `autodev_install`, and `autodev_audit` tools directly to Cursor, Claude Desktop, and other LLM agents.
- **Security Dependency Auditing**: Deep repository auditing against the OSV vulnerability database (`autodev audit`).
- **Reproducible Environments**: Exporting environment config states into reproducible `.autodev.lock.json` manifests.
- **Local Web Console**: Interactive browser console cockpit launched via `autodev ui`.
- **Benchmark Command**: CLI execution and token optimization metrics dashboard (`autodev benchmark`).

---

## [0.1.0] - 2026-05-29

### Added
- Initial release 🎉
- Core CLI with `autodev scan` and `autodev setup`
- Support for Node.js, Python, Go, Rust, Java, Docker detection
- HTML, JSON, and Markdown report generation
- Multi-platform binaries (Linux x86_64/arm64, macOS, Windows)

[Unreleased]: https://github.com/HEETMEHTA18/autodev/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/HEETMEHTA18/autodev/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/HEETMEHTA18/autodev/releases/tag/v0.1.0
