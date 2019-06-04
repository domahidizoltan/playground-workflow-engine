package stockprice

type StockDataService struct {
	repo StockDataRepository
}

func NewStockDataService(repository StockDataRepository) StockDataService {
	return StockDataService {
		repo: repository,
	}
}

func (service StockDataService) GetLatestBySymbol(symbol string, offset int, limit int) []StockData {
	return service.repo.GetLatestBySymbol(symbol, offset, limit)
}

func (service StockDataService) Save(data StockData) StockData {
	return service.repo.Save(data)
}
