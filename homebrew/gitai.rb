# GitAI Homebrew Formula
# This file should be placed in a separate tap repository
# Repository name: homebrew-tap
# File location: Formula/gitai.rb

class Gitai < Formula
  desc "AI-powered Git commit message generator using local Ollama models"
  homepage "https://github.com/xyue92/gitai"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-darwin-arm64"
      sha256 "PUT_SHA256_HERE_FOR_ARM64"
    else
      url "https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-darwin-amd64"
      sha256 "PUT_SHA256_HERE_FOR_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-linux-arm64"
      sha256 "PUT_SHA256_HERE_FOR_LINUX_ARM64"
    else
      url "https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-linux-amd64"
      sha256 "PUT_SHA256_HERE_FOR_LINUX_AMD64"
    end
  end

  depends_on "ollama"

  def install
    bin.install Dir["gitai-*"].first => "gitai"
  end

  test do
    assert_match "gitai", shell_output("#{bin}/gitai --help")
  end

  def caveats
    <<~EOS
      GitAI has been installed!

      Before using GitAI, make sure to:
      1. Start Ollama: ollama serve
      2. Pull an AI model: ollama pull qwen2.5-coder:7b

      Get started:
        cd your-git-repository
        gitai commit

      Configuration:
        gitai config --init
    EOS
  end
end
