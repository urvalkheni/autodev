# AutoDev Roadmap

## ✅ v0.1.0 — Foundation
- [x] Monorepo structure (Go workspace + PNPM workspaces)
- [x] CLI scaffold with Cobra + BubbleTea (Go CLI with TUI components)
- [x] Repository scanner (30+ languages/frameworks)
- [x] GitHub user scanner (autodev github)
- [x] Basic runtime installer (Node.js, Go, Python, Rust, etc.)
- [x] HTML / JSON / Markdown report generation (autodev report)
- [x] Neo-Brutalist marketing website (Next.js 15)
- [x] curl | bash installer script
- [x] GitHub Actions CI/CD (lint, test, build pipelines)

## ✅ v0.2.0 — Ecosystem Expansion (Current)
- [x] Docker & Kubernetes detection
  - [x] Docker CLI installer
  - [x] Detect existing Docker Desktop / run containers
  - [x] kubectl CLI installer
  - [x] Helm chart detection & installer
- [x] Mobile/SDK toolchains
  - [x] Flutter SDK installer
  - [x] Android SDK & emulator detection/installer
- [x] JVM & Build Tools
  - [x] OpenJDK installer (Java 17/21) and Kotlin support
  - [x] Maven & Gradle installers
- [x] Web & Backend languages
  - [x] PHP CLI & modules installer
  - [x] Composer installer
  - [x] Ruby runtime installer
  - [x] Bundler installer
- [x] Skills.sh integration
  - [x] Local JSON profile store (~/.config/autodev/skills-profile.json)
  - [x] Live sync with Skills.sh API (fetch curated learning paths)
- [x] Interactive Terminal UI
  - [x] BubbleTea-based skill selector (TUI)
  - [x] Visual learning roadmap viewer
- [x] Packaging & Distribution
  - [x] Homebrew tap + Scoop manifest + Chocolatey package setup
- [x] Environment Snapshot (autodev export)
  - [x] CLI entrypoint and JSON lockfile output
  - [x] Reproducible lockfile logic (captures runtimes & versions)

## 🚧 v0.3.0 — Intelligence Layer (Next)
- [ ] **AI Assistant**: Add a chat-based local LLM assistant (via Ollama/llama.cpp or similar) for interactive dev help (e.g. answer dependency/setup questions)
- [ ] **DevContainer/Cloud IDE Gen**: Auto-generate `.devcontainer.json` (or Gitpod/Codespaces config) from project scan (base images, extensions)
- [ ] **Offline / Nix Support**: Offline installation mode with caching; optional Nix flake integration for reproducible envs
- [ ] **Doctor Auto-Fixes**: Enhance `autodev doctor` with `--fix` flag to automatically remediate missing dependencies or setup issues
- [ ] **.env & Config Scaffolding**: Auto-generate `.env` or config templates based on detected frameworks (e.g. populate basic ENV vars)
- [ ] **Monorepo Support**: Multi-project workspace scanning and bulk setup (support Yarn workspace, Nx, etc.)
- [ ] **VS Code Extension**: Publish official AutoDev VSCode extension (recommend extensions, launch autodev doctor from editor)
- [ ] **Plugin System**: Design a plugin/SDK architecture so community can add custom detectors or installers
- [ ] **Team Sync**: Implement `autodev.lock` manifest for team-shared env locking and `autodev unlock` command
- [ ] **Marketplace UI (CLI)**: Interactive storefront browsing in TUI (`autodev store` / `autodev search`)

## 🌐 v1.0.0 — Production Release
- [ ] **Stable CLI API & SemVer**: Finalize CLI contract, tag v1.0.0 (long-term support)
- [ ] **Complete Documentation**: Full docs site (guides, API ref, examples) and user tutorials
- [ ] **Plugin Marketplace Launch**: Curated online marketplace and CLI integration for extensions/skills
- [ ] **Integration with Dev Platforms**: Deep integration with Codespaces/Gitpod (official config support)
- [ ] **Enterprise & Cloud Options**:
  - [ ] Enterprise self-hosted deployment (e.g. Kubernetes helm chart)
  - [ ] Single Sign-On support (OAuth/SAML) for teams
  - [ ] Hosted AutoDev SaaS with team management dashboard
- [ ] **Internationalization**: CLI and docs localization (i18n)
- [ ] **Brew/Scoop/Choco Maintenance**: Official maintained taps/buckets for easy install across OSes
- [ ] **Security & Compliance**: License scanning, vulnerability checks, and compliance certification (e.g. SOC2)
- [ ] **Final Doctor Fixes**: Ensure `autodev doctor --fix` can resolve all common setup issues automatically

## 🚀 Beyond v1.0 (Future Ideas)
- [ ] AutoDev as a Service (hosted multi-tenant cloud platform)
- [ ] Live environment dashboard (web UI to manage dev environments)
- [ ] AI-powered tool recommendations (smart suggestions based on project profile)
- [ ] First-class Swift/Nix/Docker Flake support for reproducibility
- [ ] Community-driven plugin & skill hubs (beyond the core marketplace)

## 💡 Experimental / Under Consideration
- [ ] Browser extension for one-click environment setup on GitHub repos
- [ ] Mobile app for monitoring & notifications of dev setup status
- [ ] Native NixOS integration and flake generation
- [ ] Dev Container spec generation from high-level project schema (for other IDEs)
