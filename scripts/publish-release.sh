#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
  echo "Usage: scripts/publish-release.sh <version-tag>"
  echo "Example: scripts/publish-release.sh v0.1.0"
  exit 1
fi

REPO="${BARRACUDA_REPO:-dillonlara115/barracudaseo}"
OUTPUT_DIR="dist/cli-$VERSION"

if [[ ! -d "$OUTPUT_DIR" ]]; then
  echo "Missing $OUTPUT_DIR. Run scripts/release-cli.sh $VERSION first."
  exit 1
fi

if ! command -v gh >/dev/null 2>&1; then
  echo "GitHub CLI (gh) is required. Install it and authenticate with 'gh auth login'."
  exit 1
fi

shopt -s nullglob
assets=("$OUTPUT_DIR"/*)
shopt -u nullglob

if [[ ${#assets[@]} -eq 0 ]]; then
  echo "No assets found in $OUTPUT_DIR."
  exit 1
fi

echo "Publishing $VERSION to $REPO"
if gh release view "$VERSION" --repo "$REPO" >/dev/null 2>&1; then
  echo "Release exists. Uploading assets (clobber enabled)."
  gh release upload "$VERSION" "${assets[@]}" --repo "$REPO" --clobber
else
  gh release create "$VERSION" "${assets[@]}" --repo "$REPO" --title "$VERSION" --generate-notes
fi

