package position

import (
	"time"
)

type Position struct {
	ID int					`json:"-"`
	Username string			`json:"username"`
	Symbol string			`json:"symbol"`
	Price float32			`json:"price"`
	Quantity float32		`json:"quantity"`
	StockPrice float32		`json:"stockPrice"`
	UpdatedAt time.Time		`json:"updatedAt"`
}

func createPosition(username string, symbol string) Position {
	return Position {
		Username: username,
		Symbol: symbol,
	}
}
