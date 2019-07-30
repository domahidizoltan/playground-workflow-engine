package config

import (
	"os"
	"log"
	"strings"
	"strconv"
	"github.com/spf13/viper"
)

const configPath = "config.yml"

type Server struct {
	Host, Port string
}

type Config struct {
	Redis 
	Zeebe
	App
}

type Redis struct {
	Server
}

type App struct {
	Threshold float32
}

var AppConfig Config

func init() {
	loadConfig()

	AppConfig = Config {
		Redis: getRedis(),
		Zeebe: GetZeebe(),
		App: getApp(),
	}
}

func loadConfig() {
	pwd, _ := os.Getwd()
	viper.AddConfigPath(pwd)

	viper.AutomaticEnv()
	underscoreToDot := strings.NewReplacer(".", "_",)
	viper.SetEnvKeyReplacer(underscoreToDot)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func GetServer(path string) Server {
	return Server {
		Host: viper.Get(path + ".server.host").(string),
		Port: asString(viper.Get(path + ".server.port")),
	}
}

func asString(config interface{}) string {
	var stringConfig string

	switch config.(type) {
		case int: stringConfig = strconv.Itoa(config.(int))
		default: stringConfig = config.(string)
	}

	return stringConfig
}

func getRedis() Redis {
	path := "redis"
	return Redis {
		Server: GetServer(path),
	}
}

func getApp() App {
	path := "app"
	return App {
		Threshold: float32(viper.Get(path + ".threshold").(float64)),
	}
}