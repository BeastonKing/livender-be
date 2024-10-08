package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string
	Author      string
	ReleaseYear int
	Price       int
	Image       string
	Age         int // in months
	UserID      uint
	IsSold      bool
	Genres      []*Genre `gorm:"many2many:book_genres;"`
}
