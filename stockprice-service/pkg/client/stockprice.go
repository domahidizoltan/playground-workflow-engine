package client

import (
	"strings"
	"log"
	"net/http"
	"time"
	"io"
	"io/ioutil"
	"encoding/json"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/stockprice"
)

var baseUrl string
func init() {
	baseUrl = config.AppConfig.Client.StockPriceBaseUrl
}

type jsonData struct {
	Symbol string	`json:symbol`
	Price float32	`json:price`
	Date jsonTime	`json:"updated_at,string"`
}

type StockPriceClient struct {
	service stockprice.StockDataService
}

func NewStockPriceClient(service stockprice.StockDataService) StockPriceClient{
	return StockPriceClient {
		service: service,
	}
}

func (client StockPriceClient) FetchStockData(symbol string) {
	stockSymbol := strings.ToUpper(symbol)
	
	resp, err := http.Get(baseUrl + "/" + stockSymbol)
	if err != nil {
		log.Println("Could not fetch stock data for %s : %v", stockSymbol, err)
	}
	defer resp.Body.Close()

	stockdata := parse(resp.Body)
	client.service.Save(stockdata)
}

func parse(body io.Reader) stockprice.StockData {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("Could not read response: %v", err)
	}

	data := jsonData{}
	json.Unmarshal(bytes, &data)

	stockdata := stockprice.StockData {
		Symbol: data.Symbol,
		Price: data.Price,
		Date: data.Date.time,
	}
	return stockdata
}

type jsonTime struct {
	time time.Time
}
const jsonTimeFormat = `"2006-01-02 15:04:05"`

func (j *jsonTime) UnmarshalJSON(buf []byte) error {
	parsedTime, err := time.Parse(jsonTimeFormat, string(buf))
	local, _ := time.LoadLocation("Local")
	j.time = parsedTime.In(local)
	return err
}
