package config

import (
	"io/ioutil"
	"log"
	"os"

	"fmt"
	"gopkg.in/yaml.v2"
	"reflect"
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
	URI     string   `yaml:"uri"`
	Methods []string `yaml:"methods"`
	Roles   []string `yaml:"roles"`
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

	config := ReadConfig(configFilePath)

	s := reflect.ValueOf(&config).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	/*	for key, value := range config {
		switch v := value.(type) {
		case map[string]string:
			logger.Infof("got map1 for %s", key)
		case map[string]interface{}:
			logger.Infof("got map2 for %s", key)
		case []string:
			logger.Infof("got array of string for %s", key)
		case []interface{}:
			//logger.Infof("evalauting for ", v[0])
			logger.Infof("got array of interface for %s", key)
			switch v[0].(type) {
			case map[string]string:
				logger.Infof("got array of map1 for %s", key)
			case map[string]interface{}:
				logger.Infof("got array of map2 for %s", key)
			case string:
				logger.Infof("got array of string for %s", key)
			default:
				logger.Infof("got array of unkown for %s", key)
			}
		case string:
			logger.Infof("got string value %s for %s", v, key)
		case bool:
			logger.Infof("got bool for %s", key)
		default:
			logger.Infof("got unknown type for %s", key)
		}
	}*/

	return config
}
