package trendanalyser

import "log"

type TrendAnalyser struct {
	repo SymbolHistoryRepository
}

var threshold float32

func NewTrendAnalyser(repository SymbolHistoryRepository, threshold float32) TrendAnalyser {
	threshold = threshold
	return TrendAnalyser{
		repo: repository,
	}
}

func (analyser TrendAnalyser) SetLastPrice(symbol string, price float32) {
	if err := analyser.repo.Push(symbol, price); err != nil {
		log.Println(err.Error())
	}
}

func (analyser TrendAnalyser) IsReadyToAnalyse(symbol string) bool {
	size := analyser.repo.Size(symbol)
	log.Printf("%d", size)
	return size > 0 && size % 9 == 0
}
func (analyser TrendAnalyser) Analyse(symbol string) Trend {
	lastPrices, _ := analyser.repo.Pop(symbol, 9)
	return analyser.analysePriceTrend(lastPrices)
}

func (analyser TrendAnalyser) analysePriceTrend(lastPrices []float32) Trend {
	first := analyser.mean(lastPrices[0:2])
	second := analyser.mean(lastPrices[3:5])
	third := analyser.mean(lastPrices[6:8])

	log.Printf("%f %f %f", first, second, third)
	var trend Trend
	switch {
	case first*(1+threshold) < second && second*(1+threshold) < third:
		trend = INCREASING
	case first*(1-threshold) > second && second*(1-threshold) > third:
		trend = DECREASING
	default:
		trend = STAGNATE
	}

	return trend
}

func (analyser TrendAnalyser) mean(values []float32) float32 {
	sum := float32(0.0)
	for _, value := range values {
		sum = sum + value
	}
	return sum / float32(len(values))
}
