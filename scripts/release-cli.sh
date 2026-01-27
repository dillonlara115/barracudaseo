#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
  echo "Usage: scripts/release-cli.sh <version-tag>"
  echo "Example: scripts/release-cli.sh v0.1.0"
  exit 1
fi

OUTPUT_DIR="dist/cli-$VERSION"
mkdir -p "$OUTPUT_DIR"

echo "==> Building frontend assets"
make frontend-build

echo "==> Building CLI binaries"
export CGO_ENABLED=0

build_target() {
  local goos="$1"
  local goarch="$2"
  local name="$3"
  local ext="$4"
  echo "  - $goos/$goarch -> $name$ext"
  GOOS="$goos" GOARCH="$goarch" go build -o "$OUTPUT_DIR/$name$ext" .
}

build_target darwin amd64 "barracuda_darwin_amd64" ""
build_target darwin arm64 "barracuda_darwin_arm64" ""
build_target linux amd64 "barracuda_linux_amd64" ""
build_target linux arm64 "barracuda_linux_arm64" ""
build_target windows amd64 "barracuda_windows_amd64" ".exe"

echo "==> Generating checksums"
(cd "$OUTPUT_DIR" && shasum -a 256 * > "checksums.txt")

echo
echo "Artifacts created in: $OUTPUT_DIR"
echo "Upload all files to GitHub Release: $VERSION"
echo "Tip: scripts/publish-release.sh $VERSION"
