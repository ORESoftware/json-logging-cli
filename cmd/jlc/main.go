package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/oresoftware/json-logging-cli/internal/runtimeconfig"
	logging "github.com/oresoftware/json-logging/jlog/helper"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	return runWithEnv(args, os.Environ(), stdin, stdout, stderr)
}

func runWithEnv(args []string, environ []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	cfg, err := runtimeconfig.LoadFrom(args, environ)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if cfg.Help {
		fmt.Fprint(stdout, usage())
		return 0
	}

	input := stdin
	if cfg.InputFile != "" {
		file, err := os.Open(cfg.InputFile)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		defer file.Close()
		input = file
	}

	v, err := decodeJSON(input, cfg.Strict)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	output, err := renderOutput(v, cfg)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if cfg.Format == runtimeconfig.FormatPretty && !shouldColor(cfg.Color, stdout) {
		output = stripANSI(output)
	}

	fmt.Fprintln(stdout, output)
	return 0
}

func decodeJSON(input io.Reader, strict bool) (any, error) {
	decoder := json.NewDecoder(input)
	decoder.UseNumber()

	var v any
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}
	if !strict {
		return v, nil
	}

	var extra any
	err := decoder.Decode(&extra)
	if errors.Is(err, io.EOF) {
		return v, nil
	}
	if err == nil {
		return nil, errors.New("input contains multiple JSON values")
	}
	return nil, err
}

func renderOutput(v any, cfg runtimeconfig.RuntimeConfig) (string, error) {
	switch cfg.Format {
	case runtimeconfig.FormatPretty:
		return logging.GetPrettyString(v, cfg.PrettySize), nil
	case runtimeconfig.FormatJSON:
		out, err := json.MarshalIndent(v, "", "  ")
		return string(out), err
	case runtimeconfig.FormatCompact:
		out, err := json.Marshal(v)
		return string(out), err
	default:
		return "", fmt.Errorf("unsupported output format: %s", cfg.Format)
	}
}

func shouldColor(mode runtimeconfig.ColorMode, stdout io.Writer) bool {
	switch mode {
	case runtimeconfig.ColorAlways:
		return true
	case runtimeconfig.ColorNever:
		return false
	default:
		return isTerminal(stdout)
	}
}

func isTerminal(stdout io.Writer) bool {
	file, ok := stdout.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}

func usage() string {
	return `Usage: jlc [options]

Pretty-print JSON from stdin.

Options:
  -s, --pretty-size, --size, --line-size <integer>
      Initial size passed to the json-logging pretty printer.
  -f, --input, --input-file, --file <path>
      Read JSON from this file instead of stdin.
  --format, --output-format <pretty|json|compact>
      Select output format.
  --color, --colors <auto|always|never>
      Select color behavior for pretty output.
  --strict, --strict=false
      Reject or allow trailing JSON after the first decoded value.
  -h, --help
      Print command usage and exit.

Environment:
  JLC_PRETTY_SIZE
      Default pretty size when no matching CLI flag is provided.
  JLC_INPUT_FILE
      Default input file when no matching CLI flag is provided.
  JLC_FORMAT
      Default output format: pretty, json, or compact.
  JLC_COLOR
      Default color mode: auto, always, or never.
  JLC_STRICT
      Whether to reject trailing JSON after the first decoded value.
`
}
