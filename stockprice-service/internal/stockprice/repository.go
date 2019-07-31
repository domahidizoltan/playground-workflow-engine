package stockprice

import (
	"github.com/jinzhu/gorm"
)

type StockDataRepository struct {
	db *gorm.DB
}

func NewStockDataRepository(db *gorm.DB) StockDataRepository{
	return StockDataRepository {
		db: db,
	}
}

func (repo StockDataRepository) GetLatestBySymbol(symbol string, offset int, limit int) []StockData {
	var data []StockData
	repo.db.Where("symbol = ?", symbol).
		Order("date desc").
		Offset(offset).
		Limit(limit).
		Find(&data)
	return data
}

func (repo StockDataRepository) Save(data StockData) StockData {
	repo.db.Create(&data)
	return data
}
