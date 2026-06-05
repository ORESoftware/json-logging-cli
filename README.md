# json-logging-cli

Pretty-print JSON from stdin using [`github.com/oresoftware/json-logging`](https://github.com/ORESoftware/json-logging).

```bash
echo '{"hello":"world","items":[1,2,3]}' | jlc
```

## Install

```bash
go install github.com/oresoftware/json-logging-cli/cmd/jlc@latest
```

## Build

```bash
go build -o jlc ./cmd/jlc
```

## Homebrew

Build from source through this repository's formula:

```bash
brew tap ORESoftware/json-logging-cli https://github.com/ORESoftware/json-logging-cli
brew install ORESoftware/json-logging-cli/json-logging-cli
```

Release binaries can also be published to `ORESoftware/homebrew-tap` with GoReleaser:

```bash
brew install --cask oresoftware/tap/json-logging-cli
```

To publish a release, create the `ORESoftware/homebrew-tap` repository, add a `HOMEBREW_TAP_GITHUB_TOKEN` secret with write access to that tap, then push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```
