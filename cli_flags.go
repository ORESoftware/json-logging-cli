package jsonloggingcli

import _ "embed"

// CLIFlagsTOML is the embedded flags-2-env configuration used by installed binaries.
//
//go:embed .cli-flags.toml
var CLIFlagsTOML string
