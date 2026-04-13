package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadFromFile reads and parses a YAML config file.
func LoadFromFile(path string) (*Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	return Parse(data)
}

// Parse parses YAML bytes into a Project config.
func Parse(data []byte) (*Project, error) {
	var cfg Project
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	return &cfg, nil
}
