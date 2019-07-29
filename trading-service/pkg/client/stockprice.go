package client

import (
	"strings"
	"log"
	"net/http"
	"time"
	"io"
	"io/ioutil"
	"encoding/json"
)

const baseUrl = "/api/company/real-time-price"

type StockData struct {
	Symbol string	`json:symbol`
	Price float32	`json:price`
	Date jsonTime	`json:"updated_at,string"`
}

type StockPriceClient struct {}

func NewStockPriceClient() StockPriceClient{
	return StockPriceClient {}
}

func (client StockPriceClient) FetchStockData(symbol string) StockData {
	stockSymbol := strings.ToUpper(symbol)
	
	resp, err := http.Get(baseUrl + "/" + stockSymbol)
	if err != nil {
		log.Println("Could not fetch stock data for %s : %v", stockSymbol, err)
	}
	defer resp.Body.Close()

	stockdata := parse(resp.Body)
	log.Println("Fetched", stockdata)
	return stockdata
}

func parse(body io.Reader) StockData {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("Could not read response: %v", err)
	}

	data := StockData{}
	json.Unmarshal(bytes, &data)

	stockdata := StockData {
		Symbol: data.Symbol,
		Price: data.Price,
		Date: data.Date,
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
