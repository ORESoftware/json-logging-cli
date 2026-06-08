package runtimeconfig

import (
	"strings"
	"testing"
)

func TestLoadUsesDefaultsWithoutEnvOrFlags(t *testing.T) {
	cfg, err := LoadFrom(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PrettySize != DefaultPrettySize {
		t.Fatalf("expected default pretty size %d, got %d", DefaultPrettySize, cfg.PrettySize)
	}
	if cfg.Help {
		t.Fatal("did not expect help mode")
	}
	if cfg.InputFile != "" {
		t.Fatalf("expected no input file, got %q", cfg.InputFile)
	}
	if cfg.Format != FormatPretty {
		t.Fatalf("expected default format %q, got %q", FormatPretty, cfg.Format)
	}
	if cfg.Color != ColorAuto {
		t.Fatalf("expected default color %q, got %q", ColorAuto, cfg.Color)
	}
	if !cfg.Strict {
		t.Fatal("expected strict mode by default")
	}
}

func TestLoadUsesEnvironmentValue(t *testing.T) {
	cfg, err := LoadFrom(nil, []string{EnvPrettySize + "=12"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PrettySize != 12 {
		t.Fatalf("expected env pretty size 12, got %d", cfg.PrettySize)
	}
}

func TestLoadLetsCLIOverrideEnvironment(t *testing.T) {
	cfg, err := LoadFrom([]string{
		"--pretty-size=4",
		"--input=cli.json",
		"--format=compact",
		"--color=never",
		"--strict=false",
	}, []string{
		EnvPrettySize + "=12",
		EnvInputFile + "=env.json",
		EnvFormat + "=json",
		EnvColor + "=always",
		EnvStrict + "=true",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PrettySize != 4 {
		t.Fatalf("expected CLI pretty size 4, got %d", cfg.PrettySize)
	}
	if cfg.InputFile != "cli.json" {
		t.Fatalf("expected CLI input file, got %q", cfg.InputFile)
	}
	if cfg.Format != FormatCompact {
		t.Fatalf("expected CLI format %q, got %q", FormatCompact, cfg.Format)
	}
	if cfg.Color != ColorNever {
		t.Fatalf("expected CLI color %q, got %q", ColorNever, cfg.Color)
	}
	if cfg.Strict {
		t.Fatal("expected CLI strict=false to override env")
	}
}

func TestLoadSupportsShortFlags(t *testing.T) {
	cfg, err := LoadFrom([]string{"-s5", "-h"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PrettySize != 5 {
		t.Fatalf("expected short flag pretty size 5, got %d", cfg.PrettySize)
	}
	if !cfg.Help {
		t.Fatal("expected help mode")
	}
}

func TestLoadRejectsInvalidEnvironmentInteger(t *testing.T) {
	_, err := LoadFrom(nil, []string{EnvPrettySize + "=nope"})
	if err == nil {
		t.Fatal("expected invalid env integer error")
	}
	if !strings.Contains(err.Error(), EnvPrettySize) {
		t.Fatalf("expected error to mention %s, got %v", EnvPrettySize, err)
	}
}

func TestLoadRejectsInvalidCLIInteger(t *testing.T) {
	_, err := LoadFrom([]string{"--pretty-size=nope"}, nil)
	if err == nil {
		t.Fatal("expected invalid CLI integer error")
	}
	if !strings.Contains(err.Error(), "integer") {
		t.Fatalf("expected integer parse error, got %v", err)
	}
}

func TestLoadRejectsUnknownCLIOption(t *testing.T) {
	_, err := LoadFrom([]string{"--wat"}, nil)
	if err == nil {
		t.Fatal("expected unknown option error")
	}
	if !strings.Contains(err.Error(), "--wat") {
		t.Fatalf("expected error to mention unknown option, got %v", err)
	}
}

func TestLoadRejectsPositionals(t *testing.T) {
	_, err := LoadFrom([]string{"input.json"}, nil)
	if err == nil {
		t.Fatal("expected positional argument error")
	}
	if !strings.Contains(err.Error(), "input.json") {
		t.Fatalf("expected error to mention positional argument, got %v", err)
	}
}

func TestLoadRejectsInvalidFormat(t *testing.T) {
	_, err := LoadFrom([]string{"--format=yaml"}, nil)
	if err == nil {
		t.Fatal("expected invalid format error")
	}
	if !strings.Contains(err.Error(), EnvFormat) {
		t.Fatalf("expected error to mention %s, got %v", EnvFormat, err)
	}
}

func TestLoadRejectsInvalidColorMode(t *testing.T) {
	_, err := LoadFrom([]string{"--color=sometimes"}, nil)
	if err == nil {
		t.Fatal("expected invalid color mode error")
	}
	if !strings.Contains(err.Error(), EnvColor) {
		t.Fatalf("expected error to mention %s, got %v", EnvColor, err)
	}
}
