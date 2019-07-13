package server

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/pkg/api"
)

type server struct {
	handler api.StockDataHandler
}

func InitAndStart(context config.Context) {
	server := &server {
		handler: api.NewStockDataHandler(context.StockDataService),
	}

	router := server.configRoutes()
	server.start(router)
}

func (server *server) configRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/stockdatas/{symbol}", server.handler.GetLatestBySymbol).Methods("GET")
	return router
}

func (server *server) start(router http.Handler) {
	conf := config.AppConfig.Http.Server
	log.Println("HTTP Server is running on port " + conf.Port)
	if err := http.ListenAndServe(":" + conf.Port, router); err != nil {
		panic(err)
	}
}
