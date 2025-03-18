package db

import (
	"dinushc/gorutines/configs"
	"dinushc/gorutines/pkg/dsn"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) *Db {
	dsn := dsn.GetDSN(conf)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
