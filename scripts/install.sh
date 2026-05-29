#!/usr/bin/env bash
# AutoDev Installer — https://autodev.dev
# Usage: curl -fsSL https://autodev.dev/install.sh | bash
set -euo pipefail

AUTODEV_VERSION="${AUTODEV_VERSION:-latest}"
INSTALL_DIR="${AUTODEV_INSTALL_DIR:-/usr/local/bin}"
REPO="heetmehta18/autodev"
BINARY_NAME="autodev"
GITHUB_API="https://api.github.com"
RELEASES_URL="https://github.com/${REPO}/releases"

# ── Colours ──────────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
BLUE='\033[0;34m'; BOLD='\033[1m'; RESET='\033[0m'

info()    { echo -e "${BLUE}[autodev]${RESET} $*"; }
success() { echo -e "${GREEN}[autodev]${RESET} $*"; }
warn()    { echo -e "${YELLOW}[autodev]${RESET} $*"; }
error()   { echo -e "${RED}[autodev] ERROR:${RESET} $*" >&2; exit 1; }

# ── Banner ────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}"
cat <<'EOF'
  AutoDev Installer
  Clone. Scan. Install. Build. — The App Store for Developers
  https://autodev.dev
EOF
echo -e "${RESET}"

# ── OS / Arch detection ───────────────────────────────────────────────────────
detect_os() {
  case "$(uname -s)" in
    Linux*)   echo "linux" ;;
    Darwin*)  echo "darwin" ;;
    MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
    *) error "Unsupported OS: $(uname -s)" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64)  echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    armv7*)        echo "armv7" ;;
    *) error "Unsupported architecture: $(uname -m)" ;;
  esac
}

OS=$(detect_os)
ARCH=$(detect_arch)
info "Detected: ${OS}/${ARCH}"

# ── Fetch latest version ──────────────────────────────────────────────────────
if [ "${AUTODEV_VERSION}" = "latest" ]; then
  info "Fetching latest release..."
  if command -v curl &>/dev/null; then
    AUTODEV_VERSION=$(curl -fsSL "${GITHUB_API}/repos/${REPO}/releases/latest" \
      | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
  elif command -v wget &>/dev/null; then
    AUTODEV_VERSION=$(wget -qO- "${GITHUB_API}/repos/${REPO}/releases/latest" \
      | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
  else
    error "curl or wget is required"
  fi
fi

info "Installing AutoDev v${AUTODEV_VERSION}..."

# ── Download ──────────────────────────────────────────────────────────────────
ARCHIVE_NAME="autodev_${OS}_${ARCH}"
if [ "${OS}" = "windows" ]; then
  ARCHIVE_EXT="zip"
else
  ARCHIVE_EXT="tar.gz"
fi

DOWNLOAD_URL="${RELEASES_URL}/download/v${AUTODEV_VERSION}/${ARCHIVE_NAME}.${ARCHIVE_EXT}"
TMP_DIR=$(mktemp -d)
ARCHIVE_PATH="${TMP_DIR}/autodev.${ARCHIVE_EXT}"

info "Downloading from: ${DOWNLOAD_URL}"
if command -v curl &>/dev/null; then
  curl -fsSL --progress-bar "${DOWNLOAD_URL}" -o "${ARCHIVE_PATH}"
else
  wget -q --show-progress "${DOWNLOAD_URL}" -O "${ARCHIVE_PATH}"
fi

# ── Extract ───────────────────────────────────────────────────────────────────
info "Extracting..."
if [ "${ARCHIVE_EXT}" = "tar.gz" ]; then
  tar -xzf "${ARCHIVE_PATH}" -C "${TMP_DIR}"
else
  unzip -q "${ARCHIVE_PATH}" -d "${TMP_DIR}"
fi

BINARY_PATH="${TMP_DIR}/${BINARY_NAME}"
if [ ! -f "${BINARY_PATH}" ]; then
  # Try with OS in name
  BINARY_PATH="${TMP_DIR}/${BINARY_NAME}_${OS}_${ARCH}"
fi
chmod +x "${BINARY_PATH}"

# ── Install ───────────────────────────────────────────────────────────────────
info "Installing to ${INSTALL_DIR}..."
if [ -w "${INSTALL_DIR}" ]; then
  cp "${BINARY_PATH}" "${INSTALL_DIR}/${BINARY_NAME}"
else
  sudo cp "${BINARY_PATH}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

# ── Cleanup ───────────────────────────────────────────────────────────────────
rm -rf "${TMP_DIR}"

# ── Verify ────────────────────────────────────────────────────────────────────
if command -v autodev &>/dev/null; then
  VERSION_OUT=$(autodev --version 2>&1 || true)
  success "[OK] AutoDev installed: ${VERSION_OUT}"
else
  warn "Binary installed to ${INSTALL_DIR}/${BINARY_NAME}"
  warn "Make sure ${INSTALL_DIR} is in your PATH."
fi

echo ""
echo -e "${BOLD}Getting started:${RESET}"
echo -e "  ${GREEN}autodev${RESET}              — open interactive installer"
echo -e "  ${GREEN}autodev install nodejs${RESET}   — install a specific package"
echo -e "  ${GREEN}autodev profile web-dev${RESET}  — install a developer profile"
echo -e "  ${GREEN}autodev github USERNAME${RESET}  — scan a GitHub user's repos"
echo -e "  ${GREEN}autodev doctor${RESET}           — check environment health"
echo ""
echo -e "  Docs: ${BLUE}https://autodev.dev/docs${RESET}"
echo ""
