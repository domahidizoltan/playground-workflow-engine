package api

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/jinzhu/gorm"
)

func writeResponse(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Printf("Could not convert to json: %v", data)
	}
	w.Write([]byte(jsonData))
}

func writeResponseOrError(w http.ResponseWriter, data interface{}, err error) {
	if err == nil {
		jsonData, err := json.Marshal(&data)
		if err != nil {
			log.Printf("Could not convert to json: %v", data)
		}
		w.Write([]byte(jsonData))
	} else {
		status := http.StatusInternalServerError
		if gorm.IsRecordNotFoundError(err) {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
	}
}