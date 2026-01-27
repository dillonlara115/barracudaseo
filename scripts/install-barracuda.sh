#!/usr/bin/env bash
set -euo pipefail

REPO_DEFAULT="dillonlara115/barracuda"
REPO="${BARRACUDA_REPO:-$REPO_DEFAULT}"
VERSION="${BARRACUDA_VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  darwin) OS="darwin" ;;
  linux) OS="linux" ;;
  msys*|mingw*|cygwin*) OS="windows" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

ASSET="barracuda_${OS}_${ARCH}"
if [[ "$OS" == "windows" ]]; then
  ASSET="${ASSET}.exe"
fi

if [[ "$VERSION" == "latest" ]]; then
  URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
else
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${URL}"
curl -fsSL "$URL" -o "$TMP_DIR/$ASSET"

chmod +x "$TMP_DIR/$ASSET"

if [[ "$OS" == "windows" ]]; then
  echo "Binary downloaded to: $TMP_DIR/$ASSET"
  echo "Move it somewhere in your PATH."
  exit 0
fi

if [[ ! -w "$INSTALL_DIR" ]]; then
  echo "Installing to $INSTALL_DIR (requires sudo)..."
  sudo mv "$TMP_DIR/$ASSET" "$INSTALL_DIR/barracuda"
else
  mv "$TMP_DIR/$ASSET" "$INSTALL_DIR/barracuda"
fi

echo "Installed: $INSTALL_DIR/barracuda"
