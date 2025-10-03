package config

type PulseConfig struct {
	Port         int
	TestInterval int64
}

func ParseConfig() *PulseConfig {
	// TODO Actually parse a TOML config
	return &PulseConfig{
		Port:         7150,
		TestInterval: 10,
	}
}
