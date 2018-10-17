package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ClientId          string     `yaml:"client-id"`
	ClientSecret      string     `yaml:"client-secret"`
	DiscoveryUrl      string     `yaml:"discovery-url"`
	EnableDefaultDeny string     `yaml:"enable-default-deny"`
	Listen            string     `yaml:"listen"`
	SecureCookie      string     `yaml:"secure-cookie"`
	Verbose           string     `yaml:"verbose"`
	EnableLogging     string     `yaml:"enable-logging"`
	CorsOrigins       []string   `yaml:"cors-origins"`
	CorsMethods       []string   `yaml:"cors-methods"`
	Resources         []Resource `yaml:"resources"`
}

type Resource struct {
	URI string `yaml:"uri"`
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

/*func ReadConfig(filePath string) string {
	// Read YML
	log.Println("Reading YAML Configuration")
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	return string(source)
}*/

func GetControllerConfig() Config {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "config.yaml"
	}

	config := ReadConfig(configFilePath)

	return config
}
