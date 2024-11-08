package repository

import (
	"chat-app/types"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connection() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=chatapp sslmode=disable password=postgres host=localhost")
	if err != nil {
		log.Fatal("db connection error:", err)
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

func GetChannelMessages(channelId string) ([]types.MessageDto, error) {
	db := connection()
	defer db.Close()
	tx := db.MustBegin()
	rows, err := tx.Query("SELECT messages.id, messages.text, users.id, users.username, users.profile_image FROM messages INNER JOIN users ON messages.user_id = users.id WHERE channel_id = $1", channelId)
	if err != nil {
		return nil, err
	}
	var messages []types.MessageDto
	for rows.Next() {
		var msg types.MessageDto
		if err := rows.Scan(&msg.Id, &msg.Text, &msg.User.Id, &msg.User.Username, &msg.User.ProfileImage); err != nil {
			log.Fatal(err)
		}
		messages = append(messages, msg)
	}
	tx.Commit()
	return messages, nil
}
