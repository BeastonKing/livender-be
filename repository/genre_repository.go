package repository

import (
	"livender-be/model"

	"gorm.io/gorm"
)

type GenreRepo struct {
	conn *gorm.DB
}

func NewGenreRepo(conn *gorm.DB) GenreRepo {
	return GenreRepo{conn}
}

func (gr GenreRepo) Store(genre *model.Genre) error {
	err := gr.conn.Create(genre).Error
	if err != nil {
		return err
	}
	return nil
}

func (gr GenreRepo) FindAll(genres *[]model.Genre) error {
	err := gr.conn.Find(genres).Error
	if err != nil {
		return err
	}
	return nil
}

func (gr GenreRepo) FindByID(id int, genre *model.Genre) error {
	err := gr.conn.First(genre, id).Error
	if err != nil {
		return err
	}
	return nil
}
