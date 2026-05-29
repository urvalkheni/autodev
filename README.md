<h1 align="center">
  <img src="https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/apps/website/public/logo.svg" alt="AutoDev" width="80" />
  <br/>
  AutoDev
</h1>

```text
  █████╗ ██╗   ██╗████████╗ ██████╗ ██████╗ ███████╗██╗   ██╗
  ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔══██╗██╔════╝██║   ██║
  ███████║██║   ██║   ██║   ██║   ██║██║  ██║█████╗  ██║   ██║
  ██╔══██║██║   ██║   ██║   ██║   ██║██║  ██║██╔══╝  ╚██╗ ██╔╝
  ██║  ██║╚██████╔╝   ██║   ╚██████╔╝██████╔╝███████╗ ╚████╔╝ 
  ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚══════╝ ╚═════╝ ╚══════╝  ╚═══╝ 
```


<p align="center">
  <strong>Clone. Scan. Install. Build.</strong>
</p>

<p align="center">
  An open-source, cross-platform developer environment bootstrapper that automatically detects technologies, installs missing runtimes, dependencies, SDKs, and dev tools — all with a single command.
</p>

<p align="center">
  <a href="https://github.com/HEETMEHTA18/autodev/actions/workflows/ci.yml">
    <img src="https://github.com/HEETMEHTA18/autodev/actions/workflows/ci.yml/badge.svg" alt="CI" />
  </a>
  <a href="https://github.com/HEETMEHTA18/autodev/releases/latest">
    <img src="https://img.shields.io/github/v/release/HEETMEHTA18/autodev?color=brightgreen" alt="Release" />
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="License" />
  </a>
  <a href="https://github.com/HEETMEHTA18/autodev/stargazers">
    <img src="https://img.shields.io/github/stars/HEETMEHTA18/autodev?style=social" alt="Stars" />
  </a>
</p>

---

## Quick Start

```bash
# One-line install
curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash

# Or via NPX
npx autodev setup

# Or via PNPM
pnpm dlx autodev setup

# Or via Homebrew (macOS/Linux)
brew install HEETMEHTA18/tap/autodev

# Or via Scoop (Windows)
scoop install autodev

# Or via Docker
docker run --rm -v $(pwd):/workspace ghcr.io/heetmehta18/autodev setup
```

Then run in any repo:

```bash
autodev setup
```

---

## 🔍 What It Does

| Command | Description |
|---------|-------------|
| `autodev scan` | Scan current repo for languages, frameworks, package managers |
| `autodev setup` | Install all missing runtimes and dependencies |
| `autodev github <USER>` | Scan all public repos of a GitHub user |
| `autodev doctor` | Check environment health |
| `autodev report` | Generate HTML/PDF/JSON environment report |
| `autodev skills` | Show personalized learning roadmap |
| `autodev install <tool>` | Install a specific runtime or tool |
| `autodev update` | Update all managed runtimes |
| `autodev clean` | Remove cached downloads and temp files |
| `autodev export` | Export environment config as reproducible JSON |

---

## 🧠 Detection Engine

AutoDev scans your repository and detects:

**Languages:** Node.js · Python · Go · Rust · Java · Kotlin · PHP · Ruby · C/C++ · .NET · Flutter · Dart  
**Frameworks:** React · Next.js · Angular · Vue · Svelte · Laravel · Rails · Spring Boot · Gin  
**Package Managers:** npm · pnpm · yarn · bun · pip · poetry · cargo · maven · gradle · composer  
**Containers:** Docker · Docker Compose · Kubernetes · Helm  
**Databases:** PostgreSQL · MySQL · MongoDB · Redis · SQLite  
**Cloud:** Firebase · Supabase · Terraform · Nginx · Apache  

---

## 📦 Installation Methods

| Method | Command |
|--------|---------|
| Shell | `curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh \| bash` |
| NPX | `npx autodev` |
| PNPM | `pnpm dlx autodev` |
| Homebrew | `brew install HEETMEHTA18/tap/autodev` |
| Scoop | `scoop install autodev` |
| Chocolatey | `choco install autodev` |
| Docker | `docker run ghcr.io/heetmehta18/autodev` |
| Binary | [GitHub Releases](https://github.com/HEETMEHTA18/autodev/releases) |

---

## 🏗️ Project Structure

```
autodev/
├── apps/
│   └── website/          # Next.js 15 marketing site + docs
├── packages/
│   ├── cli/              # Go CLI (cobra + bubbletea)
│   ├── core/             # OS/arch detection, config
│   ├── scanner/          # Repo + GitHub scanner
│   ├── installer/        # Runtime installer
│   ├── skills/           # skills.sh integration
│   └── github/           # GitHub API client
├── scripts/
│   ├── install.sh        # Curl installer
│   └── build.sh          # Build script
├── .github/
│   └── workflows/        # CI/CD pipelines
├── go.work               # Go workspace
├── pnpm-workspace.yaml   # PNPM workspaces
└── turbo.json            # Turborepo config
```

---

## 🤝 Contributing

We welcome contributions of all kinds! Please read [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

- 🐛 [Report a Bug](https://github.com/HEETMEHTA18/autodev/issues/new?template=bug_report.md)
- 💡 [Request a Feature](https://github.com/HEETMEHTA18/autodev/issues/new?template=feature_request.md)
- 📖 [Improve Docs](https://github.com/HEETMEHTA18/autodev/tree/main/apps/website)

---

## 📜 License

MIT © [AutoDev Contributors](https://github.com/HEETMEHTA18/autodev/graphs/contributors)
