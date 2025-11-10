class Dilivet < Formula
  desc "Diagnostics and vetting toolkit for ML-DSA (Dilithium-like) signatures"
  homepage "https://github.com/codethor0/dilivet"
  version "0.2.3"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/codethor0/dilivet/releases/download/v0.2.3/darwin-arm64.tar.gz"
      sha256 "b702337226221c5bc544eeb22a20f82e81b7f80d5367ee74e15922eb14e05417"
    else
      url "https://github.com/codethor0/dilivet/releases/download/v0.2.3/darwin-amd64.tar.gz"
      sha256 "a39654178aa008c1f27fc8aae8ba737fe9722c9a09b5f8b12db8f2f9bac903f0"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/codethor0/dilivet/releases/download/v0.2.3/linux-arm64.tar.gz"
      sha256 "744c7936366d4e23ee33bdd60fe13fecedb42fd86b6b40d824239bcbf5610f58"
    else
      url "https://github.com/codethor0/dilivet/releases/download/v0.2.3/linux-amd64.tar.gz"
      sha256 "ce56670ddb6ae854b2fe9d04512f62d1a095f93b47d7fcdb8346326e12c5c2c6"
    end
  end

  def install
    bin.install "dilivet"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/dilivet -version")
  end
end

