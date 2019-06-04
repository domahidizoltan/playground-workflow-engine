package config

import (
	"os"
	"log"
	"github.com/olebedev/config"
)

const (
	configFilePath = "config.yml"
	configNotFoundMessage = " not found in config file"

	serverGroup = "server."
	portPath = serverGroup + "port"
	priceURLPath = serverGroup + "priceUrl"
)

type ServerConfig struct {
	Port int
	PriceURL string
}

func GetConfig() ServerConfig {
	pwd, _ := os.Getwd()
	path := pwd + "/" + configFilePath
	conf, err := config.ParseYamlFile(path)
	if err != nil {
		log.Fatalln("Could not open " + path)
	}

	return ServerConfig {
		Port: getPort(conf),
		PriceURL: getPriceURL(conf),
	}
	
}

func getPort(conf *config.Config) int {
	port, err := conf.Int(portPath)
	if err != nil {
		log.Println(portPath + configNotFoundMessage)
	}

	return port
}

func getPriceURL(conf *config.Config) string {
	priceURL, err := conf.String(priceURLPath)
	if err != nil {
		log.Println(priceURLPath + configNotFoundMessage)
	}

	return priceURL
}