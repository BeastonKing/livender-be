package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID       uint
	BookID       uint
	PurchaseDate time.Time

	Book Book `gorm:"foreignKey:BookID"`
}
