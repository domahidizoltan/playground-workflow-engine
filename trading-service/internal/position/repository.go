package position

import (
	"time"
	"github.com/jinzhu/gorm"
)

type PositionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) PositionRepository{
	return PositionRepository {
		db: db,
	}
}

func (repo PositionRepository) FindAllBy(username string) ([]Position, error) {
	var data []Position
	err := repo.db.Where("username = ?", username).Order("symbol asc").Find(&data).Error
	return data, err
}

func (repo PositionRepository) FindBy(username string, symbol string) (Position, error) {
	var data Position
	err := repo.db.Where("username = ? and symbol = ?", username, symbol).Find(&data).Error
	return data, err
}

func (repo PositionRepository) Save(position Position) Position {
	position.UpdatedAt = time.Now()
	repo.db.Create(&position)
	return position
}

func (repo PositionRepository) Delete(username string, symbol string) {
	data := createPosition(username, symbol)
	repo.db.Where("username = ? and symbol = ?", username, symbol).Delete(&data)
}
