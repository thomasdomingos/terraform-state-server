package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Host string `yaml:"host",envconfig:"SERVER_HOST"`
		Port string `yaml:"port",envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	Registry struct {
		Path string `yaml:"path",envconfig:"REGISTRY_PATH"`
	} `yaml:"registry"`
	Database struct {
		Path string `yaml:"path",envconfig:"DATABASE_PATH"`
	} `yaml:"database"`
}

func processError(err error) {
	log.Fatal(err)
}

func readFile(configPath string, cfg *Config) {
	f, err := os.Open(configPath)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}

func InitConfig(configPath string, cfg *Config) error {
	readFile(configPath, cfg)
	readEnv(cfg)
	if _, err := os.Stat(cfg.Registry.Path); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.Registry.Path, os.ModePerm)
		return err
	}
	return nil
}
