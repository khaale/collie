package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Artifact represents information about system's artifact type and it's location
type Artifact struct {
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
	SearchPattern string `yaml:"searchPattern"`
}

// Config represents a configuration for collecting artifacts
type Config struct {
	Artifacts []Artifact `yaml:"artifacts"`
}

// NewConfig returns a new decoded from YAML Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
