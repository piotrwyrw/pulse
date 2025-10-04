package config

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const ConfigLocation = "/etc/pulsedaq/config.toml"

//go:embed default-config.toml
var defaultConfig []byte

type PulseMainConfig struct {
	Port int `toml:"port"`
}

type PulseTestingConfig struct {
	TestInterval    int64  `toml:"interval"`
	RecordsLocation string `toml:"records"`
}

type PulseConfig struct {
	Pulse   PulseMainConfig    `toml:"pulse"`
	Testing PulseTestingConfig `toml:"testing"`
}

func ParseConfig() (*PulseConfig, error) {
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
	cfg, err := ParseConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
