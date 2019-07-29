package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/position"
)

type PositionHandler struct {
	service position.PositionService
}

func NewPositionHandler(service position.PositionService) PositionHandler {
	return PositionHandler {
		service: service,
	}
}

func (handler PositionHandler) ListPositions(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	data, err := handler.service.GetPositions(username)
	writeResponseOrError(w, data, err)
}

func (handler PositionHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	symbol := mux.Vars(r)["symbol"]
	data, err := handler.service.GetPosition(username, symbol)
	writeResponseOrError(w, data, err)
}
