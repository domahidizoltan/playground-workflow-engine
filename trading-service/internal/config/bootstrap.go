package config

import (
	"fmt"
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	ac "github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/accountconfig"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/position"
)

type Context struct {
	AccountConfigService ac.AccountConfigService
	PositionService position.PositionService
}

func Bootstrap(conf Config) Context {
	db := connectDB(conf.DB)
	accountConfigRepository := ac.NewAccountConfigRepository(db)
	positionRepository := position.NewPositionRepository(db)
	accountConfigService := ac.NewAccountConfigService(accountConfigRepository)

	context := Context {
		AccountConfigService: accountConfigService,
		PositionService: position.NewPositionService(positionRepository, accountConfigService),
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