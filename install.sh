#!/bin/sh
set -e

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Resolve binary name
case "$OS" in
  darwin)
    BINARY_NAME="kibuild-mcp-darwin-${ARCH}"
    ;;
  linux)
    if [ "$ARCH" = "arm64" ]; then
      echo "Linux ARM64 detected."
      echo "Pre-built Linux ARM64 binaries are not yet provided."
      echo "Please build from source: https://github.com/priyabratasahoo21/kibuild-mcp#build-from-source"
      exit 1
    fi
    BINARY_NAME="kibuild-mcp-linux-${ARCH}"
    ;;
  mingw*|msys*|cygwin*)
    BINARY_NAME="kibuild-mcp-windows-amd64.exe"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

DOWNLOAD_URL="https://github.com/priyabratasahoo21/kibuild-mcp/releases/latest/download/${BINARY_NAME}"
INSTALL_DIR="/usr/local/bin"
TMP_FILE="${TMPDIR:-/tmp}/kibuild-mcp-download"

# ── Step 1: Binary ──────────────────────────────────────────────────────────

echo "Downloading ${BINARY_NAME}..."
if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$TMP_FILE" "$DOWNLOAD_URL"
else
  echo "Error: neither curl nor wget found. Please install one and retry."
  exit 1
fi

chmod +x "$TMP_FILE"

if [ ! -d "$INSTALL_DIR" ]; then
  echo "Creating ${INSTALL_DIR}..."
  sudo mkdir -p "$INSTALL_DIR"
fi

echo "Installing to ${INSTALL_DIR}/kibuild-mcp..."
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_FILE" "${INSTALL_DIR}/kibuild-mcp"
else
  echo "Requesting administrator privileges to install to ${INSTALL_DIR}..."
  sudo mv "$TMP_FILE" "${INSTALL_DIR}/kibuild-mcp"
fi

if [ "$OS" = "darwin" ]; then
  xattr -d com.apple.quarantine "${INSTALL_DIR}/kibuild-mcp" 2>/dev/null || true
fi

echo "✓ kibuild-mcp installed  ($(${INSTALL_DIR}/kibuild-mcp --version 2>/dev/null || echo 'version unknown'))"

# ── Step 2: Claude Code slash command ───────────────────────────────────────

COMMAND_DIR="${HOME}/.claude/commands"
mkdir -p "$COMMAND_DIR" 2>/dev/null || true

SETUP_URL="https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/.claude/commands/setup-kibuild.md"
INIT_URL="https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/.claude/commands/init-kibuild-project.md"

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$SETUP_URL" -o "${COMMAND_DIR}/setup-kibuild.md" 2>/dev/null && echo "✓ /setup-kibuild command installed" || true
  curl -fsSL "$INIT_URL"  -o "${COMMAND_DIR}/init-kibuild-project.md" 2>/dev/null && echo "✓ /init-kibuild-project command installed" || true
elif command -v wget >/dev/null 2>&1; then
  wget -qO "${COMMAND_DIR}/setup-kibuild.md"         "$SETUP_URL" 2>/dev/null && echo "✓ /setup-kibuild command installed" || true
  wget -qO "${COMMAND_DIR}/init-kibuild-project.md"  "$INIT_URL"  2>/dev/null && echo "✓ /init-kibuild-project command installed" || true
fi

# ── Step 3: Hand off to the native interactive setup ────────────────────────
# The binary writes the MCP config and verifies tools itself (one code path,
# identical on every OS). It needs a real terminal for the prompts — when this
# script is run via `curl | sh`, stdin is the pipe, so we redirect from /dev/tty.

BINARY_PATH="${INSTALL_DIR}/kibuild-mcp"

echo ""
if [ -e /dev/tty ]; then
  "$BINARY_PATH" --setup < /dev/tty || true
else
  echo "────────────────────────────────────────────────────────────────────────"
  echo ""
  echo "  Binary installed. Finish setup by running this in your terminal:"
  echo ""
  echo "    kibuild-mcp --setup"
  echo ""
  echo "  Or, in Claude Code, run:  /setup-kibuild"
  echo ""
  echo "  Docs: https://github.com/priyabratasahoo21/kibuild-mcp"
  echo ""
fi
