package config

import (
	"fmt"
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/domahidizoltan/playground-workflow-engine/stockprice-service/internal/stockprice"
)

type Context struct {
	StockDataService stockprice.StockDataService
}

func Bootstrap(conf Config) Context {
	db := connectDB(conf.DB)
	repo := stockprice.NewStockDataRepository(db)

	context := Context {
		StockDataService: stockprice.NewStockDataService(repo),
	}

	return context
}

func connectDB(config DB) *gorm.DB {
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
	config.Server.Host, config.Server.Port, config.User, config.Password, config.DbName)
	
	db, err := gorm.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Could not open DB connection %v", err)
	}

	return db
}