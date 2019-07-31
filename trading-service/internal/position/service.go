package position

import (
	"fmt"
	"log"
	"github.com/jinzhu/gorm"
	ac "github.com/domahidizoltan/playground-workflow-engine/trading-service/internal/accountconfig"
	"github.com/domahidizoltan/playground-workflow-engine/trading-service/pkg/client"
)

type PositionService struct {
	repo PositionRepository
	accountConfigService ac.AccountConfigService
}

func NewPositionService(repository PositionRepository, accountConfigService ac.AccountConfigService) PositionService {
	return PositionService {
		repo: repository,
		accountConfigService: accountConfigService,
	}
}

func (service PositionService) GetPositions(username string) ([]Position, error) {
	return service.repo.FindAllBy(username)
}

func (service PositionService) GetPosition(username string, symbol string) (Position, error) {
	return service.repo.FindBy(username, symbol)
}

func (service PositionService) Buy(username string, stockData client.StockData) (Position, error) {
	var err error
	var position Position
	symbol := stockData.Symbol
	log.Printf("Received %v to Buy for %s", stockData, username)

	if _, e := service.repo.FindBy(username, symbol); !gorm.IsRecordNotFoundError(e) {
		err = fmt.Errorf("%s already has an open position of %s", username, symbol)
	} else {
		accountConfig, e := service.getValidAccountConfig(username, symbol)
		if e == nil {
			position, err = service.buyStock(accountConfig, symbol, stockData.Price)
			log.Printf("Acquired position %v", position)
		} else {
			err = e
		}
	}

	return position, err
}

func (service PositionService) Sell(username string, stockData client.StockData) error {
	var err error
	symbol := stockData.Symbol
	log.Printf("Received %v to Sell for %s", stockData, username)

	if position, e := service.repo.FindBy(username, symbol); e != nil {
		err = fmt.Errorf("%s has no position of %s", username, symbol)
	} else {
		price := position.Quantity * stockData.Price
		service.accountConfigService.Deposit(username, price)
		service.repo.Delete(username, symbol)
		log.Printf("Closed position for %s on %s for price %f", username, symbol, price)
	}
	return err
}

func (service PositionService) getValidAccountConfig(username string, symbol string) (ac.AccountConfig, error) {
	var err error

	accountConfig, e := service.accountConfigService.GetAccountConfig(username)
	if e != nil {
		err = fmt.Errorf("%s has no configuration", username)
	} else if accountConfig.Balance <= 0 {
		err = fmt.Errorf("%s has 0 balance", username)
	}

	if limitConfig, ok := accountConfig.LimitConfig[symbol]; !ok {
		err = fmt.Errorf("%s has no configuration for stock %s", username, symbol)
	} else if limitConfig.BalanceLimit > accountConfig.Balance {
		err = fmt.Errorf("%s configured balance limit %f exceeds current balance")
	}

	return accountConfig, err
}

func (service PositionService) buyStock(accountConfig ac.AccountConfig, symbol string, stockPrice float32) (Position, error) {
	var position Position
	limitConfig := accountConfig.LimitConfig[symbol]

	if stockPrice > limitConfig.PriceLimit {
		return Position{}, fmt.Errorf("Stock price %f is higher than confiured price limit %f for symbol %s", stockPrice, limitConfig.PriceLimit, symbol)
	}

	quantity := limitConfig.BalanceLimit / stockPrice
	price := quantity * stockPrice

	position = Position {
		Username: accountConfig.Username,
		Symbol: symbol,
		Price: price,
		Quantity: quantity,
		StockPrice: stockPrice,
	}

	service.accountConfigService.Withdraw(accountConfig.Username, price)
	return service.repo.Save(position), nil
}