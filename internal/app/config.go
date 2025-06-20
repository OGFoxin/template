package app

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Server struct {
		BindPort         string `yaml:"bind_port"`
		LogLevel         string `yaml:"log_level"`
		StatisticRefresh int    `yaml:"statistic_refresh"`
	}
}

func (cfg *Config) getConfig(configPath string) (*Config, error) {
	yml, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yml, cfg)
	return cfg, nil

}

func NewConfig(configPath string) *Config {
	cfg := &Config{}
	actualConfig, err := cfg.getConfig(configPath)

	if err != nil {
		log.Fatal(err)
	}

	return actualConfig

}
