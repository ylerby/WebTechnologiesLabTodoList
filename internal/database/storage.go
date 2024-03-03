package database

import "gorm.io/gorm"

type Sql struct {
	DB *gorm.DB
}

func New() *Sql {
	return &Sql{}
}
