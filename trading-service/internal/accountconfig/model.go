package accountconfig

type LimitConfig struct {
	PriceLimit float32		`json:"priceLimit"`
	BalanceLimit float32	`json:"balanceLimit"`
}

type AccountConfig struct {
	Username string						`json:"username" gorm:"unique"`
	Balance float32						`json:"balance"`
	LimitConfig map[string]LimitConfig	`json:"limitConfig" gorm:"-"`
	LimitConfigString string			`json:"-" gorm:"column:limit_config"`
}

func createAccountConfig(username string) AccountConfig {
	return AccountConfig {
		Username: username,
	}
}

func createAccountConfigWithBalance(username string, balance float32) AccountConfig {
	return AccountConfig {
		Username: username,
		Balance: balance,
	}
}