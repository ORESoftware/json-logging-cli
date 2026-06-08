package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPrintsHelpWithoutReadingStdin(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithEnv([]string{"--help"}, nil, strings.NewReader("not json"), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "--format") {
		t.Fatalf("expected usage to mention format flag, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRunCompactFormat(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithEnv([]string{"--format=compact"}, nil, strings.NewReader(`{"b":2,"a":1}`), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); got != `{"a":1,"b":2}` {
		t.Fatalf("unexpected compact output: %q", got)
	}
}

func TestRunJSONFormat(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithEnv([]string{"--format=json"}, nil, strings.NewReader(`{"hello":"world"}`), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "\n  \"hello\": \"world\"\n") {
		t.Fatalf("unexpected JSON output: %q", stdout.String())
	}
}

func TestRunAutoColorStripsANSIForNonTerminalOutput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithEnv(nil, nil, strings.NewReader(`{"hello":"world"}`), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if strings.Contains(stdout.String(), "\x1b[") {
		t.Fatalf("expected no ANSI codes in auto color pipe output: %q", stdout.String())
	}
}

func TestRunStrictModeRejectsTrailingJSON(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithEnv(nil, nil, strings.NewReader(`{"a":1} {"b":2}`), &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "multiple JSON values") {
		t.Fatalf("expected trailing JSON error, got %q", stderr.String())
	}
}

func TestRunCanAllowTrailingJSON(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runWithEnv([]string{"--strict=false", "--format=compact"}, nil, strings.NewReader(`{"a":1} {"b":2}`), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); got != `{"a":1}` {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunReadsInputFile(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	if err := os.WriteFile(inputPath, []byte(`{"from":"file"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runWithEnv([]string{"--input", inputPath, "--format=compact"}, nil, strings.NewReader(`{"from":"stdin"}`), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); got != `{"from":"file"}` {
		t.Fatalf("unexpected file output: %q", got)
	}
}
