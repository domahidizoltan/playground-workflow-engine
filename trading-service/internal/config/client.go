package config

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
	baseUrl := "http://" + server.Host + ":" + server.Port

	return StockPrice {
		Server: server,
		BaseUrl: baseUrl,
	}
}