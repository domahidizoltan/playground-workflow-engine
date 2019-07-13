package stockprice

import (
	"fmt"
	"time"
)

type PriceData struct {
	Symbol string	   `json:"symbol"`
	Price jsonPrice	   `json:"price"`
	UpdatedAt jsonTime `json:"updated_at"`
}

type jsonPrice float32
type jsonTime time.Time
const jsonTimeFormat = `"2006-01-02 15:04:05"`

func (ct jsonTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	timeString := t.Format(jsonTimeFormat)
	return []byte(timeString), nil
}

func (price jsonPrice) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", price)), nil
}

