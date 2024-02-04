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
	Kafka struct {
		BrokerList     []string `yaml:"brokerList"`
		ConsumerTopics []string `yaml:"consumerTopics"`
		ProducerTopic  string   `yaml:"producerTopic"`
		GroupID        string   `yaml:"groupId"`
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
