package repository

import (
	"livender-be/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	conn *gorm.DB
}

func NewUserRepo(conn *gorm.DB) UserRepo {
	return UserRepo{conn}
}

func (ur UserRepo) Store(user *model.User) error {
	err := ur.conn.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepo) FindAll(users *[]model.User) error {
	err := ur.conn.Find(users).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepo) FindByID(id int, user *model.User) error {
	err := ur.conn.First(user, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepo) Update(user *model.User) error {
	err := ur.conn.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepo) Delete(user *model.User) error {
	err := ur.conn.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepo) FindByUsername(username string, user *model.User) error {
	return ur.conn.Where("username = ?", username).First(user).Error
}
