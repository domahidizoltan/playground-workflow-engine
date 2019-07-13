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
	DB 
	Http
	App 
	Client
	Zeebe
}

type DB struct {
	Server
	User, Password, DbName string
}

type App struct {
	Default
}

type Default struct {
	PageOffset, PageLimit int
}

var AppConfig Config

func init() {
	loadConfig()

	AppConfig = Config {
		DB: getDB(),
		Http: GetHttp(),
		App: getApp(),
		Client: GetClient(),
		Zeebe: GetZeebe(),
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

func getDB() DB {
	path := "db"
	return DB {
		Server: GetServer(path),
		User: viper.Get(path + ".user").(string),
		Password: viper.Get(path + ".password").(string),
		DbName: viper.Get(path + ".dbName").(string),
	}
}

func getApp() App {
	path := "app"
	return App {
		Default: Default {
			PageOffset: viper.Get(path + ".default.pageOffset").(int),
			PageLimit: viper.Get(path + ".default.pageLimit").(int),	
		},
	}
}
