#!/bin/sh
set -e

REPO="Naly-programming/devid"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux*)  OS="linux" ;;
  darwin*) OS="darwin" ;;
  mingw*|msys*|cygwin*) OS="windows" ;;
  *) echo "Unsupported OS: $OS" && exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

# Get latest version
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
if [ -z "$VERSION" ]; then
  echo "Failed to get latest version"
  exit 1
fi

echo "Installing devid v$VERSION ($OS/$ARCH)..."

# Build download URL
EXT="tar.gz"
if [ "$OS" = "windows" ]; then
  EXT="zip"
fi
URL="https://github.com/$REPO/releases/download/v$VERSION/devid_${VERSION}_${OS}_${ARCH}.${EXT}"

# Download and extract
TMP=$(mktemp -d)
trap "rm -rf $TMP" EXIT

echo "Downloading $URL"
curl -fsSL "$URL" -o "$TMP/devid.$EXT"

if [ "$EXT" = "zip" ]; then
  unzip -q "$TMP/devid.$EXT" -d "$TMP"
else
  tar -xzf "$TMP/devid.$EXT" -C "$TMP"
fi

# Install
BINARY="devid"
if [ "$OS" = "windows" ]; then
  BINARY="devid.exe"
fi

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"
else
  echo "Need sudo to install to $INSTALL_DIR"
  sudo mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"
fi

chmod +x "$INSTALL_DIR/$BINARY"

echo "devid v$VERSION installed to $INSTALL_DIR/$BINARY"
echo ""
echo "Get started:"
echo "  devid init"
