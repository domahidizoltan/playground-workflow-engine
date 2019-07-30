package trendanalyser

import (
	"time"
)


type Trend int

const (
	INCREASING Trend = 1
	STAGNATE Trend = 0
	DECREASING Trend = -1
)

type StockData struct {
	Symbol string	`json:"symbol"`
	Price float32	`json:"price"`
	Date time.Time	`json:"date"`
}
