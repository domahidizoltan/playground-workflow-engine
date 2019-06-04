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

	// x := context.StockDataService.GetLastByName("aa", 1)
	// log.Println("last %v", x)

	// data := stockprice.StockData {
	// 	Name: "aa",
	// 	Price: 12.34,
	// 	Date: time.Now(),
	// }
	// newData := context.StockDataService.Save(data)
	// log.Printf("saved %v", newData)
	return context
}

func connectDB(config DB) *gorm.DB {
	dbUrl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
	config.Host, config.Port, config.User, config.Password, config.DbName)
	
	db, err := gorm.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Could not open DB connection %v", err)
	}

	return db
}