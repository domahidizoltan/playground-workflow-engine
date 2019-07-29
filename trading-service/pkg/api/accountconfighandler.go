package api

import (
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	ac "github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/accountconfig"
)

type AccountConfigHandler struct {
	service ac.AccountConfigService
}

func NewAccountConfigHandler(service ac.AccountConfigService) AccountConfigHandler {
	return AccountConfigHandler {
		service: service,
	}
}

func (handler AccountConfigHandler) GetAccountConfig(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	data, err := handler.service.GetAccountConfig(username)
	writeResponseOrError(w, data, err)
}

func (handler AccountConfigHandler) CreateAccountConfig(w http.ResponseWriter, r *http.Request) {
	var config ac.AccountConfig
	_ = json.NewDecoder(r.Body).Decode(&config)
	newConfig, err := handler.service.CreateAccountConfig(config.Username, config.Balance)
	writeResponseOrError(w, newConfig, err)
	w.WriteHeader(http.StatusCreated)
}

func (handler AccountConfigHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	amount := getFloatParam(r, "amount", 0.0)
	handler.service.Deposit(username, amount)
}

func (handler AccountConfigHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	amount := getFloatParam(r, "amount", 0.0)
	err := handler.service.Withdraw(username, amount)
	writeResponseOrError(w, err.Error, err)
}

func (handler AccountConfigHandler) SetLimitConfig(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	var config map[string]ac.LimitConfig
	_ = json.NewDecoder(r.Body).Decode(&config)
	err := handler.service.SetLimitConfig(username, config)
	writeResponseOrError(w, config, err)
}

func getFloatParam(r *http.Request, key string, defaultValue float64) float32 {
	param := r.URL.Query().Get(key)
	value, err := strconv.ParseFloat(param, 64)
	if err != nil {
		value = defaultValue
	}
	return float32(value)
}
