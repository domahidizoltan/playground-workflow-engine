package server

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/config"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/pkg/api"
)

type server struct {
	accountConfigHandler api.AccountConfigHandler
	positionHandler api.PositionHandler
}

func InitAndStart(context config.Context) {
	server := &server {
		accountConfigHandler: api.NewAccountConfigHandler(context.AccountConfigService),
		positionHandler: api.NewPositionHandler(context.PositionService),
	}

	router := server.configRoutes()
	server.start(router)
}

func (server *server) configRoutes() http.Handler {
	router := mux.NewRouter()
	server.configureAccountRouts(router)
	server.configurePositionRouts(router)
	return router
}

func (server *server) configureAccountRouts(router *mux.Router) {
	const path = "/api/accounts"
	router.HandleFunc(path, server.accountConfigHandler.CreateAccountConfig).Methods("POST")
	router.HandleFunc(path + "/{username}", server.accountConfigHandler.GetAccountConfig).Methods("GET")
	router.HandleFunc(path + "/{username}/limit", server.accountConfigHandler.SetLimitConfig).Methods("PUT")
	router.HandleFunc(path + "/{username}/deposit", server.accountConfigHandler.Deposit).Methods("POST")
	router.HandleFunc(path + "/{username}/withdraw", server.accountConfigHandler.Withdraw).Methods("POST")
}

func (server *server) configurePositionRouts(router *mux.Router) {
	const path = "/api/positions/{username}"
	router.HandleFunc(path, server.positionHandler.ListPositions).Methods("GET")
	router.HandleFunc(path + "/{symbol}", server.positionHandler.GetPosition).Methods("GET")
}

func (server *server) start(router http.Handler) {
	conf := config.AppConfig.Http.Server
	log.Println("HTTP Server is running on port " + conf.Port)
	if err := http.ListenAndServe(":" + conf.Port, router); err != nil {
		panic(err)
	}
}
