package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	Port           string `yaml:"port"`
	MySQLDSN       string `yaml:"mysql_dsn"`
	SessionSecret  string `yaml:"session_secret"`
	CORSOrigin     string `yaml:"cors_origin"`
	DevSeed        bool   `yaml:"dev_seed"`
	EnableDevLogin bool   `yaml:"enable_dev_login"`
}

const defaultPath = "config.yaml"

// Load reads the YAML config file. Path resolution order:
//  1. CODEZONE_CONFIG env var, if set
//  2. ./config.yaml relative to the process working directory
//
// The file is required — there is no in-code default set. Missing or invalid
// files cause the caller to fail fast.
func Load() (Config, error) {
	path := os.Getenv("CODEZONE_CONFIG")
	if path == "" {
		path = defaultPath
	}
	return LoadFrom(path)
}

func LoadFrom(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("config: read %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("config: parse %s: %w", path, err)
	}
	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("config: %s: %w", path, err)
	}
	return cfg, nil
}

func (c Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}
	if c.MySQLDSN == "" {
		return fmt.Errorf("mysql_dsn is required")
	}
	if c.SessionSecret == "" {
		return fmt.Errorf("session_secret is required")
	}
	if c.CORSOrigin == "" {
		return fmt.Errorf("cors_origin is required")
	}
	return nil
}
