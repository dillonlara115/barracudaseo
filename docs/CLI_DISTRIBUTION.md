# CLI Distribution (Releases + Homebrew)

This repo supports distributing the CLI as prebuilt binaries via GitHub Releases, with an optional Homebrew tap pointing to those binaries.

## 1) Build release artifacts

```bash
scripts/release-cli.sh v0.1.0
```

This produces binaries + `checksums.txt` in `dist/cli-v0.1.0/`.

Upload all files to the matching GitHub Release tag (e.g. `v0.1.0`).
You can use the helper script (requires the GitHub CLI):

```bash
scripts/publish-release.sh v0.1.0
```

## 2) Install script (non-Homebrew users)

Host the install script or point users to:

```bash
curl -fsSL https://raw.githubusercontent.com/dillonlara115/barracudaseo/main/scripts/install-barracuda.sh | bash
```

Optional overrides:
- `BARRACUDA_VERSION=v0.1.0`
- `BARRACUDA_REPO=org/repo`
- `INSTALL_DIR=$HOME/.local/bin`

## 3) Homebrew tap (binary formula)

Create a separate tap repo, e.g. `barracuda/homebrew-tap`, with a formula like:

```ruby
class Barracuda < Formula
  desc "Barracuda CLI"
  homepage "https://github.com/dillonlara115/barracudaseo"
  version "0.3.0"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/dillonlara115/barracudaseo/releases/download/v0.1.0/barracuda_darwin_arm64"
    sha256 "<sha256>"
  elsif OS.mac?
    url "https://github.com/dillonlara115/barracudaseo/releases/download/v0.1.0/barracuda_darwin_amd64"
    sha256 "<sha256>"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/dillonlara115/barracudaseo/releases/download/v0.1.0/barracuda_linux_arm64"
    sha256 "<sha256>"
  elsif OS.linux?
    url "https://github.com/dillonlara115/barracudaseo/releases/download/v0.1.0/barracuda_linux_amd64"
    sha256 "<sha256>"
  end

  def install
    bin.install Dir["barracuda_*"].first => "barracuda"
  end
end
```

Then users install with:

```bash
brew install barracuda/tap/barracuda
```

## Notes about private repos

If the repo is private, GitHub Releases are also private. You’ll need:
- a public mirror repo for binaries, or
- a private installer that uses a GitHub token, or
- a storage bucket (S3/GCS) for release assets.
