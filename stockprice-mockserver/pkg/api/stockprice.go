package api

import (
	"encoding/json"
	"strings"
	"log"
	"strconv"
	"net/http"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-mockserver/internal/stockprice"
)

type MockServer struct {
	Generator *stockprice.PriceGenerator
}

func NewMockServer() *MockServer {
	return &MockServer{
		Generator: stockprice.NewGenerator(),
	}
}

func (mockServer *MockServer) GetPriceHandle(resp http.ResponseWriter, req *http.Request) {
	trendValue := req.URL.Query().Get("trend")
	newTrend := mockServer.parseTrend(trendValue)
	symbols := strings.Split(req.URL.Path, "/")
	symbol := ""
	if len(symbols) >= 4 {
		symbol = symbols[4]
	}

	priceData := mockServer.Generator.GenerateFor(symbol, newTrend)
	message, err := json.Marshal(*priceData)
	if err != nil {
		log.Fatalln("Could not marshall price data for symbol " + symbol)
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(message))
}

func (mockServer *MockServer) parseTrend(trendValue string) float32 {
	newTrend := mockServer.Generator.Trend
	parsedTrend, err := strconv.ParseFloat(trendValue, 32)
	if trendValue != "" {
		if err != nil {
			log.Println("Could not parse trend value " + trendValue)
		}

		newTrend = float32(parsedTrend)	
		if newTrend != mockServer.Generator.Trend {
			mockServer.Generator.Trend = newTrend
			log.Printf("Changed trend to %.2f", newTrend)
		}
	
	} 

	return newTrend
}

