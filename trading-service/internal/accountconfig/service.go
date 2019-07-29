package accountconfig

import (
	"fmt"
	"log"
	"strings"
	"github.com/jinzhu/gorm"
)

type AccountConfigService struct {
	repo AccountConfigRepository
}

func NewAccountConfigService(repository AccountConfigRepository) AccountConfigService {
	return AccountConfigService {
		repo: repository,
	}
}

func (service AccountConfigService) GetAccountConfig(username string) (AccountConfig, error) {
	return service.repo.FindBy(username)
}

func (service AccountConfigService) CreateAccountConfig(username string, balance float32) (AccountConfig, error) {
	return service.repo.Save(username, balance)
}

func (service AccountConfigService) Deposit(username string, amount float32) {
	if _, err := service.repo.FindBy(username); err != nil {
		service.repo.Save(username, amount)
	} else {
		service.repo.UpdateBalance(username, amount)
	}
	log.Printf("Deposited %f for %s", amount, username) 
}

func (service AccountConfigService) Withdraw(username string, amount float32) error {
	var err error
	config, e := service.repo.FindBy(username)
	if e != nil {
		err = fmt.Errorf("%s has no AccountConfig", username)
	} else if config.Balance - amount < 0 {
		err = fmt.Errorf("Could not withdraw %f from AccountConfig %f for %s", amount, config.Balance, username)
	} else {
		service.repo.UpdateBalance(username, -amount)
		log.Printf("Withdrawn %f for %s", amount, username)
	}

	return err
}

func (service AccountConfigService) SetLimitConfig(username string, inputConfig map[string]LimitConfig) error {
	var err error
	config, e := service.repo.FindBy(username)
	if gorm.IsRecordNotFoundError(e) == true {
		err = fmt.Errorf("%s has no AccountConfig", username)
	} else if limitConfig, e := service.sanitizeConfig(inputConfig, config.Balance); e == nil {
		service.repo.UpdateLimitConfig(username, limitConfig)
	} else {
		err = e
	}

	return err
}

func (service AccountConfigService) sanitizeConfig(inputConfig map[string]LimitConfig, balance float32) (map[string]LimitConfig, error) {
	var err error
	limitConfig := make(map[string]LimitConfig, len(inputConfig))

	sum := float32(0.0)
	for symbol, conf := range inputConfig {
		sanitizedConfig := LimitConfig {
			PriceLimit: conf.PriceLimit,
			BalanceLimit: conf.BalanceLimit,
		}

		if (sanitizedConfig.PriceLimit <= 0) {
			err = fmt.Errorf("Price limit for %s must be a positive number", symbol)
		}

		sum = sum + conf.PriceLimit

		limitConfig[strings.ToUpper(symbol)] = sanitizedConfig
	}

	if sum > balance {
		err = fmt.Errorf("Config price sum %f exceeds balance %f", sum, balance)
	}

	return limitConfig, err
}