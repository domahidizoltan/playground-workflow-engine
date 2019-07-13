package stockprice

import (
	"time"
)

type StockData struct {
	ID int			`json:"-"`
	Symbol string	`json:symbol`
	Price float32	`json:price`
	Date time.Time	`json:date`
}
