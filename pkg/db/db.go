package db

import (
	"dinushc/gorutines/configs"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) *Db {
	dsn := GetDSN(conf)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}

func GetDSN(conf *configs.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Db.Host, conf.Db.Port, conf.Db.User, conf.Db.Password, conf.Db.Name)
}
