package configs

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Database struct {
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Db   string `yaml:"db"`
	} `yaml:"database"`
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
}

func ReadConfig() (*Config, error) {
	configFile, err := ioutil.ReadFile("src/configs/local.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
