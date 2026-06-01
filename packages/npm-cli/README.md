# AutoDev CLI — The App Store for Developers ⚡

[![npm version](https://img.shields.io/npm/v/@heetmehta18/autodev.svg)](https://www.npmjs.com/package/@heetmehta18/autodev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AutoDev** is an open-source, cross-platform developer environment bootstrapper. It acts as an **App Store for Developers**, simplifying complex toolchain setups through intelligent, profile-based automation.

---

## 🚀 Key Features

*   🔍 **Polyglot Codebase Scanner**: Detects 30+ languages, frameworks, package managers, and DevOps infrastructure.
*   🛡️ **Supply-Chain Security Audits**: Queries the OSV (Open Source Vulnerabilities) database to find safety risks in your dependencies.
*   ⚙️ **Ballast Installer**: Automatically extracts and links compilers/runtimes (like Node.js, Go, Python, Rust) to your path.
*   📦 **Monorepo / Multi-project Scanner**: Groups and maps nested modules dynamically within a monorepo structure.
*   🛰️ **Cloud IDE Scaffolding**: Runs `autodev containerize` to generate `.devcontainer.json` environment setups.
*   🔄 **Config Migrator**: Seamlessly upgrades legacy profile configurations to standard YAML.

---

## 📦 Installation & Quick Start

You can run AutoDev on the fly using Node's package executor, or install it globally.

### 1. Run Instantly (No Installation)
Scan your workspace and bootstrap dependencies without installing anything permanently:
```bash
npx @heetmehta18/autodev setup
```

### 2. Install Globally
Install the package globally for instant local terminal access:
```bash
npm install -g @heetmehta18/autodev
```

### 3. Usage
Verify that the CLI is installed and ready:
```bash
autodev --help
```

---

## 🛠️ Main CLI Commands

### 🔍 scan
Analyzes your current working directory for configuration markers, lockfiles, and monorepo folders:
```bash
autodev scan
```

### 📦 setup
Scans the project and aligns your local development environment by downloading missing runtimes:
```bash
autodev setup
autodev setup --yes  # Skip confirmation prompts
```

### 🩺 doctor
Inspects your system configuration, checks tool versions against the `.autodev.lock.json` lockfile, and scans for exposed secrets (like AWS keys or GitHub tokens):
```bash
autodev doctor
autodev doctor --fix  # Restore lockfile mismatches automatically
```

### 🛡️ audit
Scans lockfiles and dependencies for known supply-chain vulnerabilities using the OSV database:
```bash
autodev audit
```

### 💻 containerize
Generates dev container setup configurations (`.devcontainer.json`) and VSCode plugin recommendations based on the detected stack:
```bash
autodev containerize
```

### 🔄 migrate
Upgrades legacy `.json` profile configs to the standard `.autodev.yaml` schema:
```bash
autodev migrate
```

---

## ⚙️ How It Works

1.  **Platform Detection**: The wrapper maps your OS and CPU architecture (e.g. `linux/amd64`, `darwin/arm64`) to the corresponding compiled Go binary.
2.  **Binary Caching**: Downloads the pre-compiled binary directly from the matching GitHub Release tag. Subsequent executions are run instantly from cache.
3.  **Process Delegation**: Delegated execution forwards all streams, signals, exit codes, and Model Context Protocol (MCP) servers seamlessly.
