package accountconfig

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type AccountConfigRepository struct {
	db *gorm.DB
}

func NewAccountConfigRepository(db *gorm.DB) AccountConfigRepository{
	return AccountConfigRepository {
		db: db,
	}
}

func (repo AccountConfigRepository) FindBy(username string) (AccountConfig, error) {
	var data AccountConfig
	err := repo.db.Table("account_configs").
			Where("username = ?", username).
			Scan(&data).Error

	if len(data.LimitConfigString) > 0 {
		err = json.Unmarshal([]byte(data.LimitConfigString), &data.LimitConfig)
	}

	return data, err
}

func (repo AccountConfigRepository) Save(username string, balance float32) (AccountConfig, error) {
	data := createAccountConfigWithBalance(username, balance)
	err := repo.db.Create(&data).Error
	return data, err
}

func (repo AccountConfigRepository) UpdateBalance(username string, balance float32) {
	repo.db.Table("account_configs").
		Where("username = ?", username).
		Update("balance", gorm.Expr("balance + ?", balance))
}

func (repo AccountConfigRepository) UpdateLimitConfig(username string, config map[string]LimitConfig) {
	limitConfig, _ := json.Marshal(config)
	repo.db.Table("account_configs").
		Where("username = ?", username).
		Update("limit_config", limitConfig)
}