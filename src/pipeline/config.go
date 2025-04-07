package pipeline

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the pipeline configuration
type Config struct {
	Name        string   `yaml:"name"`
	Repository  string   `yaml:"repository"`
	Branch      string   `yaml:"branch"`
	BuildCmd    string   `yaml:"build_cmd"`
	TestCmd     string   `yaml:"test_cmd"`
	DeployCmd   string   `yaml:"deploy_cmd"`
	Artifacts   []string `yaml:"artifacts"`
	Environment map[string]string `yaml:"environment"`
}

// LoadConfig loads the pipeline configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}