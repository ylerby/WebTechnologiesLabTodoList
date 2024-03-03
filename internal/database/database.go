package database

import (
	"backend/internal/model"
)

func (d *Database) CreateUser(login, password string) (bool, error) {
	var alreadyExists model.UserModel

	queryResult := d.DB.Where("login = ?", login).Find(&alreadyExists)
	if err := queryResult.Error; err != nil {
		return false, err
	}

	if queryResult.RowsAffected != 0 {
		return true, nil
	}

	user := &model.UserModel{
		Login:    login,
		Password: password,
	}

	queryResult = d.DB.Create(user)

	if err := queryResult.Error; err != nil {
		return false, err
	}

	return false, nil
}

func (d *Database) GetUser(login string) (*model.UserModel, bool, error) {
	var result model.UserModel

	queryResult := d.DB.Where("login = ?", login).Find(&result)
	if err := queryResult.Error; err != nil {
		return nil, true, err
	}

	if queryResult.RowsAffected == 0 {
		return nil, false, nil
	}

	return &result, true, nil
}
