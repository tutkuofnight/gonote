package main

import (
	database "chat-app/db"
	"chat-app/middleware"
	"chat-app/routes"
	"chat-app/types"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New()
	db := database.GetDB()
	app.Static("/", "./static")
	var clients = make(map[*websocket.Conn]bool)
	channelUsers := make(map[string]int32)

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:name", middleware.RestrictUser, websocket.New(func(conn *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)

		user := conn.Locals("user").(map[string]interface{})
		userId := int(user["id"].(float64))

		var channel types.Channel
		channelName := conn.Params("name")
		channelErr := db.Model(&types.Channel{}).Where("name = ?", channelName).Preload("Messages").First(&channel).Error
		if channelErr != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Channel is not found"))
			conn.Close()
			return
		}
		currentCount, has := channelUsers[channelName]
		if has {
			channelUsers[channelName] = currentCount + 1
		} else {
			channelUsers[channelName] = 1
		}
		defer func() {
			channelUsers[channelName] -= 1
			fmt.Printf("%s kanalina bagli kullanici sayisi: %d\n", channelName, channelUsers[channelName])
		}()
		fmt.Printf("%s kanalina bagli kullanici sayisi: %d\n", channelName, channelUsers[channelName])
		clients[conn] = true

		for {
			var message types.Message
			if mt, msg, err = conn.ReadMessage(); err != nil {
				log.Info("read:", err)
				break
			}
			json.Unmarshal(msg, &message)
			message.ChannelId = channel.Id
			message.UserId = int(userId)
			err := db.Create(&message).Error
			if err != nil {
				log.Info("create error:", err)
			}

			var user types.User
			if err := db.Find(&user, message.UserId).Association("Channels").Append(&channel); err != nil {
				log.Error("update error:", err)
			}
			for client := range clients {
				if err = client.WriteMessage(mt, msg); err != nil {
					log.Info("write:", err)
					break
				}
			}
		}
	}))

	routes.ChannelRoutes(app)
	routes.UserRotues(app)

	log.Fatal(app.Listen(":3000"))

}
