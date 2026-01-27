class Barracuda < Formula
  desc "Barracuda CLI"
  homepage "https://github.com/dillonlara115/barracuda"
  version "0.0.0"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/dillonlara115/barracuda/releases/download/v0.0.0/barracuda_darwin_arm64"
    sha256 "REPLACE_WITH_SHA256"
  elsif OS.mac?
    url "https://github.com/dillonlara115/barracuda/releases/download/v0.0.0/barracuda_darwin_amd64"
    sha256 "REPLACE_WITH_SHA256"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/dillonlara115/barracuda/releases/download/v0.0.0/barracuda_linux_arm64"
    sha256 "REPLACE_WITH_SHA256"
  elsif OS.linux?
    url "https://github.com/dillonlara115/barracuda/releases/download/v0.0.0/barracuda_linux_amd64"
    sha256 "REPLACE_WITH_SHA256"
  else
    odie "Unsupported platform"
  end

  def install
    bin.install Dir["barracuda_*"].first => "barracuda"
  end
end
