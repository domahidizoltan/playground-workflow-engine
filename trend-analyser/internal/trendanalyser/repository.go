package trendanalyser

import (
	"strconv"
	"github.com/go-redis/redis"
)

type SymbolHistoryRepository struct {
	client *redis.Client
}

func NewSymbolHistoryRepository(client *redis.Client) SymbolHistoryRepository {
	return SymbolHistoryRepository{
		client: client,
	}
}

func (repo SymbolHistoryRepository) Push(symbol string, price float32) error {
	return repo.client.RPush(symbol, price).Err()
}

func (repo SymbolHistoryRepository) Size(symbol string) int {
	return int(repo.client.LLen(symbol).Val())
}

func (repo SymbolHistoryRepository) Pop(symbol string, size int) ([]float32, error) {
	prices := make([]float32, size)
	results, err := repo.client.LRange(symbol, 0, int64(size-1)).Result()
	if err == nil {
		for i, result := range results {
			price, _ := strconv.ParseFloat(result, 32)
			prices[i] = float32(price)
		}
		repo.client.LTrim(symbol, int64(size), 0)
	}

	return prices, err
}
