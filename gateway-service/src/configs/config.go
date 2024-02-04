package configs

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Kafka struct {
		BrokerList []string `yaml:"brokerList"`
		Topic string `yaml:"topic"`
	} `yaml:"kafka"`
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
