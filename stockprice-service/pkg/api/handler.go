package api

import (
	"strings"
	"strconv"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/stockprice"
)

type StockDataHandler struct {
	service stockprice.StockDataService
}

func NewStockDataHandler(service stockprice.StockDataService) StockDataHandler {
	return StockDataHandler {
		service: service,
	}
}

func (handler StockDataHandler) GetLatestBySymbol(w http.ResponseWriter, r *http.Request) {
	offset := getIntParam(r, "offset", config.AppConfig.App.Default.PageOffset)
	limit := getIntParam(r, "limit", config.AppConfig.App.Default.PageLimit)
	symbol := mux.Vars(r)["symbol"]

	data := []stockprice.StockData{}
	if symbol != "" {
		data = handler.service.GetLatestBySymbol(strings.ToUpper(symbol), offset, limit)
	}

	writeResponse(w, data)
}

func getIntParam(r *http.Request, key string, defaultValue int) int {
	param := r.URL.Query().Get(key)
	value, err := strconv.Atoi(param)
	if err != nil {
		value = defaultValue
	}
	return value
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Printf("Could not convert to json: %v", data)
	}
	w.Write([]byte(jsonData))
}