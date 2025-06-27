package app

import (
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Server struct {
		BindPort         string `yaml:"bind_port"`
		LogLevel         string `yaml:"log_level"`
		StatisticRefresh int    `yaml:"statistic_refresh"`
		UseGin           bool   `yaml:"use_gin"`
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

func (cfg *Config) Watchdog(config string) (*Config, error) {
	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()

	if err = watcher.Add(config); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				return cfg.getConfig(config)
			}
		case err := <-watcher.Errors:
			return nil, err
		}
	}

}
