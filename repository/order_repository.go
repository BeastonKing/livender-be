package repository

import (
	"livender-be/model"

	"gorm.io/gorm"
)

type OrderRepo struct {
	conn *gorm.DB
}

func NewOrderRepo(conn *gorm.DB) OrderRepo {
	return OrderRepo{conn}
}

func (or OrderRepo) Store(order *model.Order) error {
	err := or.conn.Create(order).Error
	if err != nil {
		return err
	}
	return nil
}

func (or OrderRepo) FindByID(id int, order *model.Order) error {

	err := or.conn.Where("id = ?", id).Preload("Book").First(order).Error
	if err != nil {
		return err
	}
	return nil
}

func (or OrderRepo) FindByBookID(bookID uint) (*model.Order, error) {
	var order model.Order
	err := or.conn.Where("book_id = ?", bookID).First(&order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (or OrderRepo) FindAllByUserID(userID int) ([]model.Order, error) {
	var orders []model.Order
	err := or.conn.Where("user_id = ?", userID).Preload("Book").Preload("Book.Genres").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
