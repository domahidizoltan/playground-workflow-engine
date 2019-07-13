package stockprice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"time"
	"testing"
)

const (
	symbol = "TEST"
	baseTime = "2019-04-28 13:27:"
)

var (
	baseTick = createTick("00")
	firstTick = createTick("01")
	secondTick = createTick("02")
	thirdTick = createTick("03")
)

func setUp() {
	testTime := baseTick
	rand.Seed(1)
	now = func() time.Time { 
		testTime = testTime.Add(time.Second)
		return testTime 	
	}
}

func TestPriceGenerator_GenerateForSymbol(test *testing.T) {
	setUp()

	expectedData := []*PriceData {
		makePriceData(1.0004377, firstTick),
		makePriceData(1.0002811, secondTick),
		makePriceData(1.0010949, thirdTick),
	}

	generator := &PriceGenerator{Price: 1, Trend: 0}

	actualData := []*PriceData{
		generator.GenerateFor(symbol, 0),
		generator.GenerateFor(symbol, 0),
		generator.GenerateFor(symbol, 0),
	}

	for i, expected := range expectedData {
		assertSameSymbol(test, expected, actualData[i])
		assertSamePrice(test, expected, actualData[i])
		assertSameDate(test, expected, actualData[i])
	}

}

func TestPriceGenerator_GenerateForSymbolWithPositiveTrend(test *testing.T) {
	setUp()

	expectedData := []*PriceData {
		makePriceData(1.0004377, firstTick),
		makePriceData(2.8629165, secondTick),
		makePriceData(3.8351545, thirdTick),
	}

	generator := &PriceGenerator{Price: 1, Trend: 0}

	actualData := []*PriceData{
		generator.GenerateFor(symbol, 0),
		generator.GenerateFor(symbol, 2),
		generator.GenerateFor(symbol, 1),
	}

	for i, expected := range expectedData {
		assertSameSymbol(test, expected, actualData[i])
		assertSamePrice(test, expected, actualData[i])
		assertSameDate(test, expected, actualData[i])
	}

}

func TestPriceGenerator_GenerateForSymbolWithNegativeTrend(test *testing.T) {
	setUp()

	expectedData := []*PriceData {
		makePriceData(1.0004377, firstTick),
		makePriceData(-0.8623543, secondTick),
		makePriceData(-1.8329648, thirdTick),
	}

	generator := &PriceGenerator{Price: 1, Trend: 0}

	actualData := []*PriceData{
		generator.GenerateFor(symbol, 0),
		generator.GenerateFor(symbol, -2),
		generator.GenerateFor(symbol, -1),
	}

	for i, expected := range expectedData {
		assertSameSymbol(test, expected, actualData[i])
		assertSamePrice(test, expected, actualData[i])
		assertSameDate(test, expected, actualData[i])
	}

}

func TestPriceGenerator_GenerateForSymbolWithRestoringTrend(test *testing.T) {
	setUp()

	expectedData := []*PriceData {
		makePriceData(1.0004377, firstTick),
		makePriceData(10.313458, secondTick),
		makePriceData(10.32185, thirdTick),
	}

	generator := &PriceGenerator{Price: 1, Trend: 0}

	actualData := []*PriceData{
		generator.GenerateFor(symbol, 0),
		generator.GenerateFor(symbol, 10),
		generator.GenerateFor(symbol, 0),
	}

	for i, expected := range expectedData {
		assertSameSymbol(test, expected, actualData[i])
		assertSamePrice(test, expected, actualData[i])
		assertSameDate(test, expected, actualData[i])
	}

}

func createTick(seconds string) time.Time {
	tick, err := time.Parse(jsonTimeFormat, fmt.Sprintf("\"%s%s\"", baseTime, seconds))
	if err != nil {
		fmt.Printf("Could not create tick: %+v", err)
	}

	return tick
}

func makePriceData(price float32, tick time.Time) *PriceData {
	return &PriceData{
		Symbol: symbol,
		Price: jsonPrice(price),
		UpdatedAt: jsonTime(tick),
	}
}

var assertSameSymbol = func (test *testing.T, expected *PriceData, actual *PriceData) { 
	assert.Equal(test, expected.Symbol, actual.Symbol, "symbol differs") 
}
var assertSamePrice = func (test *testing.T, expected *PriceData, actual *PriceData) { 
	assert.Equal(test, expected.Price, actual.Price, "price differs") 
}
var assertSameDate = func (test *testing.T, expected *PriceData, actual *PriceData) { 
	assert.Equal(test, format(expected.UpdatedAt), format(actual.UpdatedAt), "date differs") 
}
var format = func (t jsonTime) string { return time.Time(t).Format(jsonTimeFormat) }
