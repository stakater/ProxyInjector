package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	GatekeeperImage   string   `yaml:"gatekeeper-image"`
	ClientId          string   `yaml:"client-id"`
	ClientSecret      string   `yaml:"client-secret"`
	DiscoveryUrl      string   `yaml:"discovery-url"`
	EnableDefaultDeny string   `yaml:"enable-default-deny"`
	Listen            string   `yaml:"listen"`
	SecureCookie      string   `yaml:"secure-cookie"`
	Verbose           string   `yaml:"verbose"`
	EnableLogging     string   `yaml:"enable-logging"`
	CorsOrigins       []string `yaml:"cors-origins"`
	CorsMethods       []string `yaml:"cors-methods"`
	Resources         []struct {
		URI     string   `yaml:"uri"`
		Methods []string `yaml:"methods"`
		Roles   []string `yaml:"roles"`
	} `yaml:"resources"`
	Scopes []string `yaml:"scopes"`
}

func ReadConfig(filePath string) Config {
	var config Config
	// Read YML
	log.Println("Reading YAML Configuration")
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	// Unmarshall
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Panic(err)
	}

	return config
}

func GetControllerConfig() Config {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "config.yaml"
	}

	return ReadConfig(configFilePath)
}
