package database

import (
	. "api_example/app/types"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	dsn := "postgres://postgres:postgres@localhost:5432/todoapp"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	if err := db.AutoMigrate(Todo{}); err != nil {
		fmt.Println(err)
	}
	if err := db.AutoMigrate(User{}); err != nil {
		fmt.Println(err)
	}
}

func GetConnection() *gorm.DB {
	return db
}
