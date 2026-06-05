class JsonLoggingCli < Formula
  desc "Pretty-print JSON from stdin with ORESoftware json-logging"
  homepage "https://github.com/ORESoftware/json-logging-cli"
  url "https://github.com/ORESoftware/json-logging-cli.git", branch: "main"
  version "0.1.0"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"jlc", "./cmd/jlc"
  end

  test do
    output = pipe_output("#{bin}/jlc", "{\"hello\":\"world\"}")
    assert_match "hello", output
    assert_match "world", output
  end
end
