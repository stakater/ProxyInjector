package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func ReadConfig(filePath string) map[string]string {
	var config map[string]string
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

func GetControllerConfig() map[string]string {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "config.yaml"
	}

	config := ReadConfig(configFilePath)

	return config
}
