package stockprice

import (
	"math/rand"
	"time"
	"strings"
)

type PriceGenerator struct {
	Price float32
	Trend float32
}

var now func() time.Time

func NewGenerator() *PriceGenerator {
	now = time.Now
	rand.Seed(now().UnixNano())
	return &PriceGenerator{
		Price: rand.Float32() * 1000,
		Trend: 0,
	}
}

func (generator *PriceGenerator) GenerateFor(symbol string, trend float32) *PriceData {
	generator.Trend = trend
	generator.Price = randomPrice(generator.Price, generator.Trend)

	return &PriceData{
		Symbol: strings.ToUpper(symbol),
		Price: jsonPrice(generator.Price),
		UpdatedAt: jsonTime(now()),
	}
}

func randomPrice(price float32, trend float32) float32 {
	change := trend + randomDeviation(trend)
	newPrice := price + randomDeviation(price * 0.01) + change
	return newPrice
}

func randomDeviation(value float32) float32 {
	sign := 1
	if rand.Float32() <= 0.5 {
		sign = -1
	}

	return rand.Float32() * (value * 0.1) * float32(sign)
}