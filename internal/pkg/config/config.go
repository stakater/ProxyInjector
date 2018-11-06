package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
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

func GetControllerConfig() []string {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "config.yaml"
	}

	config := ReadConfig(configFilePath)

	return structToStringArray(config)
}

func structToStringArray(config Config) []string {

	configArgs := []string{}

	if config.ClientId != "" {
		configArgs = append(configArgs, "--client-id="+config.ClientId)
	}
	if config.ClientSecret != "" {
		configArgs = append(configArgs, "--client-secret="+config.ClientSecret)
	}
	if config.DiscoveryUrl != "" {
		configArgs = append(configArgs, "--discovery-url="+config.DiscoveryUrl)
	}
	/*	if config.EnableDefaultDeny !="" {
		configArgs = append(configArgs, " --enable-default-deny=\""+config.EnableDefaultDeny+"\"")
	}*/
	if config.Listen != "" {
		configArgs = append(configArgs, "--listen="+config.Listen)
	}
	if config.SecureCookie != "" {
		configArgs = append(configArgs, "--secure-cookie="+config.SecureCookie)
	}
	if config.Verbose != "" {
		configArgs = append(configArgs, "--verbose="+config.Verbose)
	}
	if config.EnableLogging != "" {
		configArgs = append(configArgs, "--enable-logging="+config.EnableLogging)
	}
	for _, origin := range config.CorsOrigins {
		configArgs = append(configArgs, "--cors-origins="+origin)
	}
	for _, method := range config.CorsMethods {
		configArgs = append(configArgs, "--cors-methods="+method)
	}
	for _, resource := range config.Resources {
		//  --resources "uri=/admin*|roles=admin,superuser|methods=POST,DELETE
		res := ""
		if resource.URI != "" {
			res = "uri=" + resource.URI
		}
		if len(resource.Methods) != 0 {
			res = res + "|methods=" + strings.Join(resource.Methods, ",")
		}
		configArgs = append(configArgs, "--resources="+res)
	}

	return configArgs
}
