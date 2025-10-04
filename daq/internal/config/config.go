package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const ConfigLocation = "/etc/pulsedaq/config.toml"

//go:embed default-config.toml
var defaultConfig []byte

type PulseConfig struct {
	Http struct {
		Port int `toml:"port"`
	} `toml:"http"`
	Testing struct {
		Interval    int64  `toml:"interval"`
		RecordsPath string `toml:"records_path"`
		BinaryPath  string `toml:"binary"`
	} `toml:"speedtest"`
}

func validateConfig(cfg *PulseConfig) error {
	if cfg.Testing.RecordsPath == "" {
		return fmt.Errorf("invalid config: missing records location")
	}

	if cfg.Testing.Interval < 1 {
		return fmt.Errorf("invalid config: test interval must be >= 1")
	}

	return nil
}

func parseConfig() (*PulseConfig, error) {
	f, err := os.Open(ConfigLocation)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg PulseConfig
	_, err = toml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func createDefaultConfig() error {
	err := os.MkdirAll(filepath.Dir(ConfigLocation), 0755)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(ConfigLocation, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(defaultConfig)
	if err != nil {
		return err
	}

	return nil
}

func Initialize() (*PulseConfig, error) {
	_, err := os.Stat(ConfigLocation)
	if err != nil && os.IsNotExist(err) {
		err := createDefaultConfig()
		if err != nil {
			return nil, err
		}
		goto parse
	}
	if err != nil {
		return nil, err
	}

parse:
	cfg, err := parseConfig()
	if err != nil {
		return nil, err
	}

	err = validateConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
