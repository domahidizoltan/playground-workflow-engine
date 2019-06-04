package server

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/pkg/api"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/pkg/client"
)

type server struct {
	handler api.StockDataHandler
}

func InitAndStart() {
	context := config.Bootstrap(config.AppConfig)
	server := &server {
		handler: api.NewStockDataHandler(context.StockDataService),
	}


	spclient := client.NewStockPriceClient(context.StockDataService)
	spclient.FetchStockData("bb")


	router := server.configRoutes()
	server.start(router)
}

func (server *server) configRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/stockdatas/{symbol}", server.handler.GetLatestBySymbol).Methods("GET")
	return router
}

func (server *server) start(router http.Handler) {
	conf := config.AppConfig
	log.Println("Server is running on port " + conf.Server.Port)
	if err := http.ListenAndServe(":" + conf.Server.Port, router); err != nil {
		panic(err)
	}
}
