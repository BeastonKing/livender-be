package repository

import (
	"livender-be/model"

	"gorm.io/gorm"
)

type BookRepo struct {
	conn *gorm.DB
}

func NewBookRepo(conn *gorm.DB) BookRepo {
	return BookRepo{conn}
}

func (br BookRepo) Store(book *model.Book) error {
	err := br.conn.Create(book).Error
	if err != nil {
		return err
	}
	return nil
}

func (br BookRepo) FindAll(books *[]model.Book) error {
	err := br.conn.Preload("Genres").Find(books).Error
	if err != nil {
		return err
	}
	return nil
}

func (br BookRepo) FindByID(id int, book *model.Book) error {
	err := br.conn.Preload("Genres").Where("id = ?", id).First(book).Error
	if err != nil {
		return err
	}
	return nil
}

func (br BookRepo) Update(book *model.Book) error {
	err := br.conn.Save(book).Error
	if err != nil {
		return err
	}
	return nil
}

func (br BookRepo) Delete(book *model.Book) error {
	err := br.conn.Delete(book).Error
	if err != nil {
		return err
	}
	return nil
}

func (br BookRepo) ClearGenres(book *model.Book) error {
	return br.conn.Model(book).Association("Genres").Clear()
}

func (br BookRepo) FindAllBooksOwnedByUser(userId int, books *[]model.Book) error {
	err := br.conn.Preload("Genres").Where("user_id = ?", userId).Find(books).Error
	if err != nil {
		return err
	}
	return nil
}
