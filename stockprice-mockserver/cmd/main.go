package main

import (
	"log"
	"strconv"
	"net/http"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-mockserver/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-mockserver/pkg/api"
)

func main() {
	mockServer := api.NewMockServer()
	conf := config.GetConfig()
	http.HandleFunc(conf.PriceURL, mockServer.GetPriceHandle)

	port := strconv.Itoa(conf.Port)
	log.Println("Server is running on port " + port)
	if err:=http.ListenAndServe(":" + port, nil); err != nil {
		panic(err)
	}
}
