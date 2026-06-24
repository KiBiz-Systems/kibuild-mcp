#!/bin/sh
set -e

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
fi

BINARY_NAME="kibuild-mcp-${OS}-${ARCH}"
# If windows/cygwin/msys
if [ "${OS#*mingw}" != "$OS" ] || [ "${OS#*msys}" != "$OS" ]; then
  OS="windows"
  BINARY_NAME="kibuild-mcp-windows-amd64.exe"
fi

LATEST_RELEASE_URL="https://github.com/priyabratasahoo21/kibuild-mcp/releases/latest/download/${BINARY_NAME}"
INSTALL_DIR="/usr/local/bin"

echo "Downloading ${BINARY_NAME}..."
curl -fsSL "$LATEST_RELEASE_URL" -o kibuild-mcp

echo "Installing to ${INSTALL_DIR}/kibuild-mcp..."
chmod +x kibuild-mcp

# Try moving with sudo if needed
if [ -w "$INSTALL_DIR" ]; then
  mv kibuild-mcp "${INSTALL_DIR}/kibuild-mcp"
else
  echo "Requesting administrator privileges to install to ${INSTALL_DIR}..."
  sudo mv kibuild-mcp "${INSTALL_DIR}/kibuild-mcp"
fi

echo "Success! kibuild-mcp installed to ${INSTALL_DIR}/kibuild-mcp"
