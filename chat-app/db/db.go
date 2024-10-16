package db

import (
	. "chat-app/types"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	dsn := "postgres://postgres:postgres@localhost:5432/chatapp"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	if err := db.AutoMigrate(User{}); err != nil {
		fmt.Println(err)
	}
	if err := db.AutoMigrate(Message{}); err != nil {
		fmt.Println(err)
	}
	if err := db.AutoMigrate(Channel{}); err != nil {
		fmt.Println(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
