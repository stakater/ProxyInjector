package config

import (
	"io/ioutil"
	"log"
	"os"
)

func ReadConfig(filePath string) string {
	// Read YML
	log.Println("Reading YAML Configuration")
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	return string(source)
}

func GetControllerConfig() string {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "config.yaml"
	}

	config := ReadConfig(configFilePath)

	return config
}
