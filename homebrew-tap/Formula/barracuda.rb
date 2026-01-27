class Barracuda < Formula
  desc "Barracuda CLI"
  homepage "https://github.com/dillonlara115/barracudaseo"
  version "0.1.0"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/dillonlara115/barracudaseo/releases/download/v0.1.0/barracuda_darwin_arm64"
    sha256 "3b00592033a73865c8a3aa6ad8d1cff88b0351b9bda32a6dcfad726008564ca1"
  elsif OS.mac?
    url "https://github.com/dillonlara115/barracudaseo/releases/download/v0.1.0/barracuda_darwin_amd64"
    sha256 "113200564a2625862c255646d7b77590318b7b6f25fe6057821ace017a5bb13f"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/dillonlara115/barracudaseo/releases/download/v0.1.0/barracuda_linux_arm64"
    sha256 "b183e15f94c58c2a6de70bff2ea84673b82295f024229daf4e3cf4dfb0235615"
  elsif OS.linux?
    url "https://github.com/dillonlara115/barracudaseo/releases/download/v0.1.0/barracuda_linux_amd64"
    sha256 "c28f2f18306aeee949a3f1ac9be3b1e09507350e2efcc3812c05379b235e1866"
  else
    odie "Unsupported platform"
  end

  def install
    bin.install Dir["barracuda_*"].first => "barracuda"
  end
end
