package database

import (
	"backend/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type dbInterface interface {
	CreateUser(login, password string) (bool, error)
	GetUser(login string) (*model.UserModel, bool, error)
}

type Database struct {
	DB *gorm.DB
}

func New() *Database {
	return &Database{}
}

func (d *Database) Connect() error {
	var err error
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable"

	d.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("ошибка при подключении к БД - %s", err)
	}

	err = d.DB.AutoMigrate(&model.UserModel{})
	if err != nil {
		return fmt.Errorf("ошибка при миграции - %s", err)
	}
	return nil
}
