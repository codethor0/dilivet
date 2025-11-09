class Dilivet < Formula
  desc "Diagnostics and vetting toolkit for ML-DSA (Dilithium-like) signatures"
  homepage "https://github.com/codethor0/dilivet"
  version "0.1.11"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/codethor0/dilivet/releases/download/v0.1.11/dilivet_darwin_arm64.tar.gz"
      sha256 "TODO_REPLACE_WITH_ACTUAL_SHA256"
    else
      url "https://github.com/codethor0/dilivet/releases/download/v0.1.11/dilivet_darwin_amd64.tar.gz"
      sha256 "TODO_REPLACE_WITH_ACTUAL_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/codethor0/dilivet/releases/download/v0.1.11/dilivet_linux_arm64.tar.gz"
      sha256 "TODO_REPLACE_WITH_ACTUAL_SHA256"
    else
      url "https://github.com/codethor0/dilivet/releases/download/v0.1.11/dilivet_linux_amd64.tar.gz"
      sha256 "TODO_REPLACE_WITH_ACTUAL_SHA256"
    end
  end

  def install
    bin.install "dilivet"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/dilivet -version")
  end
end

