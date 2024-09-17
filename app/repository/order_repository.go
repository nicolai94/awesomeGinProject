package repository

import (
	"awesomeProject/app/domain/dao"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *dao.Order) (*dao.Order, error)
}

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func (u OrderRepositoryImpl) CreateOrder(order *dao.Order) (*dao.Order, error) {

	err := u.db.Create(order).Error
	if err != nil {
		log.Error("Got an error creating order. Error: ", err)
		return nil, err
	}

	return order, nil
}

func OrderRepositoryInit(db *gorm.DB) *OrderRepositoryImpl {
	db.AutoMigrate(&dao.Order{})
	return &OrderRepositoryImpl{
		db: db,
	}
}
