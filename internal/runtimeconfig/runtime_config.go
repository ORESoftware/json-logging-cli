package runtimeconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	jsonloggingcli "github.com/oresoftware/json-logging-cli"
)

const (
	EnvPrettySize     = "JLC_PRETTY_SIZE"
	EnvInputFile      = "JLC_INPUT_FILE"
	EnvFormat         = "JLC_FORMAT"
	EnvColor          = "JLC_COLOR"
	EnvStrict         = "JLC_STRICT"
	EnvHelp           = "JLC_HELP"
	EnvParseErrors    = "JLC_PARSE_ERRORS"
	EnvUnknownOptions = "JLC_UNKNOWN_OPTIONS"
	EnvPositionals    = "JLC_POSITIONALS"
)

const DefaultPrettySize = 0

type OutputFormat string

const (
	FormatPretty  OutputFormat = "pretty"
	FormatJSON    OutputFormat = "json"
	FormatCompact OutputFormat = "compact"
)

type ColorMode string

const (
	ColorAuto   ColorMode = "auto"
	ColorAlways ColorMode = "always"
	ColorNever  ColorMode = "never"
)

type RuntimeConfig struct {
	PrettySize int
	InputFile  string
	Format     OutputFormat
	Color      ColorMode
	Strict     bool
	Help       bool
}

type ParseError struct {
	Issues []string
}

func (e ParseError) Error() string {
	return strings.Join(e.Issues, "; ")
}

func Load(args []string) (RuntimeConfig, error) {
	return LoadFrom(args, os.Environ())
}

func LoadFrom(args []string, environ []string) (RuntimeConfig, error) {
	env := EnvMap(environ)
	cli, err := parseCLI(args, jsonloggingcli.CLIFlagsTOML)
	if err != nil {
		return RuntimeConfig{}, err
	}
	if err := validateCLI(cli); err != nil {
		return RuntimeConfig{}, err
	}

	return FromMap(MergeMaps(env, cli))
}

func EnvMap(environ []string) map[string]string {
	env := make(map[string]string, len(environ))
	for _, entry := range environ {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			env[key] = value
		}
	}
	return env
}

func MergeMaps(env map[string]string, cli map[string]string) map[string]string {
	combined := make(map[string]string, len(env)+len(cli))
	for key, value := range env {
		combined[key] = value
	}
	for key, value := range cli {
		combined[key] = value
	}
	return combined
}

func FromMap(values map[string]string) (RuntimeConfig, error) {
	prettySize, err := intValue(values, EnvPrettySize, DefaultPrettySize)
	if err != nil {
		return RuntimeConfig{}, err
	}
	if prettySize < 0 {
		return RuntimeConfig{}, fmt.Errorf("%s must be greater than or equal to 0", EnvPrettySize)
	}

	format, err := outputFormatValue(values, EnvFormat, FormatPretty)
	if err != nil {
		return RuntimeConfig{}, err
	}

	color, err := colorModeValue(values, EnvColor, ColorAuto)
	if err != nil {
		return RuntimeConfig{}, err
	}

	strict, err := boolValue(values, EnvStrict, true)
	if err != nil {
		return RuntimeConfig{}, err
	}

	help, err := boolValue(values, EnvHelp, false)
	if err != nil {
		return RuntimeConfig{}, err
	}

	return RuntimeConfig{
		PrettySize: prettySize,
		InputFile:  strings.TrimSpace(values[EnvInputFile]),
		Format:     format,
		Color:      color,
		Strict:     strict,
		Help:       help,
	}, nil
}

func validateCLI(cli map[string]string) error {
	var issues []string

	for _, value := range stringList(cli[EnvParseErrors]) {
		issues = append(issues, value)
	}
	for _, value := range stringList(cli[EnvUnknownOptions]) {
		issues = append(issues, "unknown option: "+value)
	}
	for _, value := range stringList(cli[EnvPositionals]) {
		issues = append(issues, "unexpected positional argument: "+value)
	}

	if len(issues) > 0 {
		return ParseError{Issues: issues}
	}
	return nil
}

func intValue(values map[string]string, key string, fallback int) (int, error) {
	value, ok := values[key]
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return parsed, nil
}

func boolValue(values map[string]string, key string, fallback bool) (bool, error) {
	value, ok := values[key]
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	parsed, ok := normalizeBool(value)
	if !ok {
		return false, fmt.Errorf("%s must be a boolean", key)
	}
	return parsed, nil
}

func outputFormatValue(values map[string]string, key string, fallback OutputFormat) (OutputFormat, error) {
	value := strings.ToLower(strings.TrimSpace(values[key]))
	if value == "" {
		return fallback, nil
	}
	switch OutputFormat(value) {
	case FormatPretty, FormatJSON, FormatCompact:
		return OutputFormat(value), nil
	default:
		return "", fmt.Errorf("%s must be one of: %s, %s, %s", key, FormatPretty, FormatJSON, FormatCompact)
	}
}

func colorModeValue(values map[string]string, key string, fallback ColorMode) (ColorMode, error) {
	value := strings.ToLower(strings.TrimSpace(values[key]))
	if value == "" {
		return fallback, nil
	}
	switch ColorMode(value) {
	case ColorAuto, ColorAlways, ColorNever:
		return ColorMode(value), nil
	default:
		return "", fmt.Errorf("%s must be one of: %s, %s, %s", key, ColorAuto, ColorAlways, ColorNever)
	}
}

func normalizeBool(value string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "t", "1", "yes", "y", "on":
		return true, true
	case "false", "f", "0", "no", "n", "off":
		return false, true
	default:
		return false, false
	}
}

func stringList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	var values []string
	if err := json.Unmarshal([]byte(value), &values); err != nil {
		return []string{value}
	}
	return values
}
