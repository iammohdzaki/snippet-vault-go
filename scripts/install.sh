#!/usr/bin/env bash
set -e

REPO="iammohdzaki/snippet-vault-go"
INSTALL_DIR="$HOME/.local/bin"

echo "Installing Snippet Vault..."

# Detect OS and Arch
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)  OS_NAME="linux" ;;
    Darwin) OS_NAME="darwin" ;;
    *)      echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64)  ARCH_NAME="amd64" ;;
    aarch64|arm64) ARCH_NAME="arm64" ;;
    *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release version
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$LATEST_TAG" ]; then
    echo "Failed to fetch latest release from GitHub."
    exit 1
fi

echo "Found latest version: $LATEST_TAG"

# Download archive
TMP_DIR=$(mktemp -d)
TAR_FILE="snippet-vault_${OS_NAME}_${ARCH_NAME}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$TAR_FILE"

echo "Downloading $DOWNLOAD_URL..."
curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$TAR_FILE"

# Extract and install
mkdir -p "$INSTALL_DIR"
tar -xzf "$TMP_DIR/$TAR_FILE" -C "$TMP_DIR"
mv "$TMP_DIR/snippet-vault" "$TMP_DIR/snippet-vault-server" "$INSTALL_DIR/"

rm -rf "$TMP_DIR"

echo "Snippet Vault installed successfully to $INSTALL_DIR!"
echo "Make sure $INSTALL_DIR is in your PATH."
