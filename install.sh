#!/usr/bin/env bash

set -e

REPO="zeflq/dockpoint"
BINARY="dockpoint"
INSTALL_DIR="/usr/local/bin"

echo "📦 Installing $BINARY from GitHub releases..."

# Detect OS
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize arch naming
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Resolve latest version
LATEST=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep tag_name | cut -d '"' -f4)

if [ -z "$LATEST" ]; then
  echo "❌ Could not fetch latest version. Check internet or GitHub status."
  exit 1
fi

echo "🔍 Detected: OS=$OS ARCH=$ARCH → version=$LATEST"

TARBALL="${BINARY}_${LATEST#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST/$TARBALL"

echo "⬇️  Downloading $TARBALL..."
curl -sSL "$URL" -o "$TARBALL"

echo "📂 Extracting..."
tar -xzf "$TARBALL"

echo "🚀 Installing to $INSTALL_DIR..."
chmod +x $BINARY
sudo mv $BINARY $INSTALL_DIR/$BINARY

echo "🧹 Cleaning up..."
rm "$TARBALL"

echo "✅ Done! Run '$BINARY --help' to get started."
