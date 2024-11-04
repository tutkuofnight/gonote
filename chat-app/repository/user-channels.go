package repository

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)


func connection() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=chatapp sslmode=disable password=postgres host=localhost")
	if err != nil {
		log.Fatal("db connection error:" , err)
	}
	return db
}
// include user to channel with many2many relations

func AddUserChannels(userId, channelId int) error {
	db := connection()
	defer db.Close()
	tx := db.MustBegin()
	_, err := tx.Exec("INSERT INTO user_channels (user_id,channel_id) VALUES ($1,$2) ON CONFLICT DO NOTHING", userId, channelId)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}