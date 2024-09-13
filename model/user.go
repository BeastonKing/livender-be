package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fullname string
	Username string `gorm:"uniqueIndex"`
	Password string
	Books    []Book

	Orders []Order
}
