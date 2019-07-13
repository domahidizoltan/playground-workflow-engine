package config

import (
	"github.com/spf13/viper"
)
type Client struct {
	StockPrice
}

type StockPrice struct {
	Server
	BaseUrl string
}

const stockPricePath = "client.stockprice"

func GetClient() Client {
	return Client {
		StockPrice: getStockPrice(),
	}
}

func getStockPrice() StockPrice {
	server := GetServer(stockPricePath)
	urlPart := viper.Get(stockPricePath + ".baseUrl").(string)
	baseUrl := "http://" + server.Host + ":" + server.Port + urlPart

	return StockPrice {
		Server: server,
		BaseUrl: baseUrl,
	}
}