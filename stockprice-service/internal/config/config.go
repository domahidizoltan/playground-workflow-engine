package config

import (
	"os"
	"log"
	"strconv"
	"strings"
	"github.com/spf13/viper"
)

const configPath = "config.yml"

type Config struct {
	Server 
	DB 
	App 
	Client client
}

type Server struct {
	Port string
}

type DB struct {
	Port int
	Host, User, Password, DbName string
}

type App struct {
	Default
}

type Default struct {
	PageOffset, PageLimit int
}

type client struct {
	StockPriceBaseUrl string
}


var AppConfig Config

func init() {
	loadConfig()

	AppConfig = Config {
		Server: getServerConfig(),
		DB: getDBConfig(),
		App: getAppConfig(),
		Client: getClientConfig(),
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

func getServerConfig() Server {
	return Server {
		Port: asString(viper.Get("server.port")),
	}
}

func getDBConfig() DB {
	return DB {
		Host: viper.Get("db.host").(string),
		Port: viper.Get("db.port").(int),
		User: viper.Get("db.user").(string),
		Password: viper.Get("db.password").(string),
		DbName: viper.Get("db.dbName").(string),
	}
}

func getAppConfig() App {
	return App {
		Default: Default {
			PageOffset: viper.Get("app.default.pageOffset").(int),
			PageLimit: viper.Get("app.default.pageLimit").(int),	
		},
	}
}

func getClientConfig() client {
	return client {
		StockPriceBaseUrl: viper.Get("client.stockPriceBaseUrl").(string),
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