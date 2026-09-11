//go:build cgo && !windows

package runtimeconfig

import (
	"os"

	flags2env "github.com/oresoftware/flags-2-env/clients/golang"
)

func parseCLI(args []string, configTOML string) (map[string]string, error) {
	file, err := os.CreateTemp("", "jlc-cli-flags-*.toml")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(configTOML); err != nil {
		file.Close()
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}

	return flags2env.ParseFromFile(file.Name(), args)
}
