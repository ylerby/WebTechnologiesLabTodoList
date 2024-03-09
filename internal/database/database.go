package database

import (
	"backend/internal/domain"
)

func (d *Database) CreateUser(login, password string) (bool, error) {
	var alreadyExists domain.UserModel

	queryResult := d.DB.Where("login = ?", login).Find(&alreadyExists)
	if err := queryResult.Error; err != nil {
		return false, err
	}

	if queryResult.RowsAffected != 0 {
		return true, nil
	}

	user := &domain.UserModel{
		Login:    login,
		Password: password,
	}

	queryResult = d.DB.Create(user)

	if err := queryResult.Error; err != nil {
		return false, err
	}

	return false, nil
}

func (d *Database) GetUser(login string) (*domain.UserModel, bool, error) {
	var result domain.UserModel

	queryResult := d.DB.Where("login = ?", login).Find(&result)
	if err := queryResult.Error; err != nil {
		return nil, true, err
	}

	if queryResult.RowsAffected == 0 {
		return nil, false, nil
	}

	return &result, true, nil
}
