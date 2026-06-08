# json-logging-cli

Pretty-print JSON from stdin using [`github.com/oresoftware/json-logging`](https://github.com/ORESoftware/json-logging).

```bash
echo '{"hello":"world","items":[1,2,3]}' | jlc
```

## Options

`jlc` uses [`ORESoftware/flags-2-env`](https://github.com/ORESoftware/flags-2-env)
with the project-local [`.cli-flags.toml`](./.cli-flags.toml) file. Runtime
configuration is loaded by merging environment variables with parsed CLI flags;
when both provide the same typed value, the CLI flag wins.

```bash
echo '{"hello":"world"}' | JLC_PRETTY_SIZE=40 jlc --pretty-size=0
echo '{"hello":"world"}' | jlc --format=json --color=never
jlc --input ./event.json --format=compact
```

Supported values:

| CLI flag | Environment | Type | Description |
| --- | --- | --- | --- |
| `-s`, `--pretty-size`, `--size`, `--line-size` | `JLC_PRETTY_SIZE` | integer | Initial size passed to the json-logging pretty printer. |
| `-f`, `--input`, `--input-file`, `--file` | `JLC_INPUT_FILE` | string | Read JSON from this file instead of stdin. |
| `--format`, `--output-format` | `JLC_FORMAT` | string | Output format: `pretty`, `json`, or `compact`. |
| `--color`, `--colors` | `JLC_COLOR` | string | Color mode: `auto`, `always`, or `never`. |
| `--strict` | `JLC_STRICT` | bool | Reject trailing JSON after the first decoded value. Defaults to `true`; pass `--strict=false` to allow trailing data. |
| `-h`, `--help` | `JLC_HELP` | bool | Print command usage and exit. |

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
