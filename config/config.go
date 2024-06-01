package config

import (
	"github.com/JensvandeWiel/logger"
	"github.com/fatih/structs"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type EnvType string

const (
	EnviromentDev  EnvType = "dev"
	EnviromentProd EnvType = "prod"
)

type Config struct {
	Environment EnvType `yaml:"environment"`
	Host        string  `yaml:"host"`
	Port        string  `yaml:"port"`
}

var conf *Config

func NewDefaultConfig() *Config {
	return &Config{
		Environment: EnviromentDev,
		Host:        "0.0.0.0",
		Port:        "8080",
	}
}

const DefaultConfigPath = "config.yaml"

func init() {
	c, err := getConfig(DefaultConfigPath)
	if err != nil {
		panic(err)
	}

	conf = c
}

func getConfig(path string) (*Config, error) {
	logger.GetLogger().Info("Reading config file", slog.String("path", path))
	config := NewDefaultConfig()

	_, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		logger.GetLogger().Warn("Config file not found, using default config")

		configBytes, err := yaml.Marshal(config)
		if err != nil {
			logger.GetLogger().Error("Failed to marshal default config")
			return nil, err
		}

		err = os.WriteFile(path, configBytes, 0644)
		if err != nil {
			logger.GetLogger().Error("Failed to write default config to file")
			return nil, err
		}

		return config, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		logger.GetLogger().Error("Failed to read config file", slog.String("path", path), slog.String("error", err.Error()))
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		logger.GetLogger().Error("Failed to unmarshal config file", slog.String("path", path), slog.String("error", err.Error()))
		return nil, err
	}

	conf = config
	return config, nil
}

func GetConfig() *Config {
	return conf
}

func GetKey(key string) (interface{}, bool) {
	val, ok := structs.Map(conf)[key]
	return val, ok
}
