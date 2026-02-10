package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type V2rayInstance struct {
	User    string `yaml:"user"`
	Passwd  string `yaml:"passwd"`
	Host    string `yaml:"host"`
	Country string `yaml:"country"`
}

var configs []V2rayInstance

func LoadConfigs(filePath string) ([]V2rayInstance, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %s", err)
	}

	for _, v := range configs {
		log.Printf("User: %s\n", v.User)
		// 不记录密码，避免泄露到日志
	}

	return configs, nil
}
