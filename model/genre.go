package model

import "gorm.io/gorm"

type Genre struct {
	gorm.Model
	Name  string
	Books []*Book `gorm:"many2many:book_genres;"`
	// Books []*Book `gorm:"many2many:book_genres;constraint:OnDelete:CASCADE;"`
}
