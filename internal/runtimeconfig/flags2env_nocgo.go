//go:build !cgo

package runtimeconfig

import (
	"encoding/json"
	"strconv"
	"strings"
)

func parseCLI(args []string, _ string) (map[string]string, error) {
	out := map[string]string{}
	var parseErrors []string
	var positionals []string
	var unknownOptions []string
	allowUnknown := false

	for i := 0; i < len(args); i++ {
		token := args[i]
		if token == "--" {
			positionals = append(positionals, args[i+1:]...)
			break
		}
		if token == "" || token == "-" || !strings.HasPrefix(token, "-") {
			positionals = append(positionals, token)
			continue
		}
		if token == "--allow-unknown" {
			allowUnknown = true
			continue
		}

		if strings.HasPrefix(token, "--") {
			i = parseLongToken(args, i, out, &parseErrors, &unknownOptions, allowUnknown)
			continue
		}

		i = parseShortToken(args, i, out, &parseErrors, &unknownOptions, allowUnknown)
	}

	setJSONList(out, EnvParseErrors, parseErrors)
	setJSONList(out, EnvPositionals, positionals)
	setJSONList(out, EnvUnknownOptions, unknownOptions)
	return out, nil
}

func parseLongToken(args []string, i int, out map[string]string, parseErrors *[]string, unknownOptions *[]string, allowUnknown bool) int {
	token := args[i]
	nameValue := strings.TrimPrefix(token, "--")
	name, value, hasValue := strings.Cut(nameValue, "=")

	switch name {
	case "pretty-size", "size", "line-size":
		if !hasValue {
			next, ok := separatedValue(args, i)
			if !ok {
				*parseErrors = append(*parseErrors, token+" requires an integer value")
				return i
			}
			value = next
			i++
		}
		if !validInteger(value) {
			*parseErrors = append(*parseErrors, token+" expects an integer value")
			return i
		}
		out[EnvPrettySize] = value
	case "input", "input-file", "file":
		if !hasValue {
			next, ok := separatedValue(args, i)
			if !ok {
				*parseErrors = append(*parseErrors, token+" requires a file path")
				return i
			}
			value = next
			i++
		}
		out[EnvInputFile] = value
	case "format", "output-format":
		if !hasValue {
			next, ok := separatedValue(args, i)
			if !ok {
				*parseErrors = append(*parseErrors, token+" requires an output format")
				return i
			}
			value = next
			i++
		}
		out[EnvFormat] = value
	case "color", "colors":
		if !hasValue {
			next, ok := separatedValue(args, i)
			if !ok {
				*parseErrors = append(*parseErrors, token+" requires a color mode")
				return i
			}
			value = next
			i++
		}
		out[EnvColor] = value
	case "strict", "single", "single-value":
		if !hasValue {
			out[EnvStrict] = "true"
			return i
		}
		parsed, ok := normalizeBool(value)
		if !ok {
			*parseErrors = append(*parseErrors, token+" expects a boolean value")
			return i
		}
		out[EnvStrict] = strconv.FormatBool(parsed)
	case "no-strict", "no-single", "no-single-value":
		out[EnvStrict] = "false"
	case "help":
		if !hasValue {
			out[EnvHelp] = "true"
			return i
		}
		parsed, ok := normalizeBool(value)
		if !ok {
			*parseErrors = append(*parseErrors, token+" expects a boolean value")
			return i
		}
		out[EnvHelp] = strconv.FormatBool(parsed)
	case "no-help":
		out[EnvHelp] = "false"
	default:
		if !allowUnknown {
			*unknownOptions = append(*unknownOptions, token)
		}
	}

	return i
}

func parseShortToken(args []string, i int, out map[string]string, parseErrors *[]string, unknownOptions *[]string, allowUnknown bool) int {
	token := args[i]
	body := strings.TrimPrefix(token, "-")
	if body == "" {
		return i
	}

	switch body[0] {
	case 's':
		value := strings.TrimPrefix(body[1:], "=")
		if value == "" {
			next, ok := separatedValue(args, i)
			if !ok {
				*parseErrors = append(*parseErrors, token+" requires an integer value")
				return i
			}
			value = next
			i++
		}
		if !validInteger(value) {
			*parseErrors = append(*parseErrors, token+" expects an integer value")
			return i
		}
		out[EnvPrettySize] = value
	case 'f':
		value := strings.TrimPrefix(body[1:], "=")
		if value == "" {
			next, ok := separatedValue(args, i)
			if !ok {
				*parseErrors = append(*parseErrors, token+" requires a file path")
				return i
			}
			value = next
			i++
		}
		out[EnvInputFile] = value
	case 'h':
		value := strings.TrimPrefix(body[1:], "=")
		if value == "" {
			out[EnvHelp] = "true"
			return i
		}
		parsed, ok := normalizeBool(value)
		if !ok {
			*parseErrors = append(*parseErrors, token+" expects a boolean value")
			return i
		}
		out[EnvHelp] = strconv.FormatBool(parsed)
	default:
		if !allowUnknown {
			*unknownOptions = append(*unknownOptions, token)
		}
	}

	return i
}

func separatedValue(args []string, i int) (string, bool) {
	if i+1 >= len(args) {
		return "", false
	}
	next := args[i+1]
	if strings.HasPrefix(next, "-") && !validInteger(next) {
		return "", false
	}
	return next, true
}

func validInteger(value string) bool {
	_, err := strconv.ParseInt(value, 10, 0)
	return err == nil
}

func setJSONList(out map[string]string, key string, values []string) {
	if len(values) == 0 {
		return
	}
	encoded, err := json.Marshal(values)
	if err == nil {
		out[key] = string(encoded)
	}
}
