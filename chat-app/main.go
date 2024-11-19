package main

import (
	"chat-app/middleware"
	"chat-app/routes"
	"chat-app/types"
	"encoding/json"
	"fmt"

	"chat-app/db"
	"chat-app/repository"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New()
	db := db.GetDB()
	app.Static("/", "./static")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000/",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

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

	app.Get("/ws/:id", middleware.RestrictUser, websocket.New(func(conn *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		user := conn.Locals("user").(map[string]interface{})
		userId := int(user["id"].(float64))
		var channel types.Channel
		channelId := conn.Params("id")
		channelErr := db.Model(&types.Channel{}).Where("id = ?", channelId).Preload("Messages").First(&channel).Error

		if channelErr != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Channel is not found"))
			conn.Close()
			return
		}

		updateOnlineCount := func() {
			onlineCount, _ := json.Marshal(channelUsers[channelId])
			for client := range clients {
				if err = client.WriteMessage(mt, onlineCount); err != nil {
					log.Info("write:", err)
					conn.Close()
					break
				}
			}
		}

		currentCount, has := channelUsers[channelId]
		if has {
			channelUsers[channelId] = currentCount + 1
			updateOnlineCount()

		} else {
			channelUsers[channelId] = 1
			updateOnlineCount()
		}
		defer func() {
			channelUsers[channelId] -= 1
			updateOnlineCount()
			fmt.Printf("%s kanalina bagli kullanici sayisi: %d\n", channelId, channelUsers[channelId])
		}()
		fmt.Printf("%s kanalina bagli kullanici sayisi: %d\n", channelId)
		clients[conn] = true

		for {
			var wsResponse types.WsResponseDto
			userMessage := types.MessageDto{}
			var rawMessage types.Message
			if mt, msg, err = conn.ReadMessage(); err != nil {
				log.Info("read:", err)
				conn.Close()
				break
			}
			jserr := json.Unmarshal(msg, &userMessage)
			if jserr != nil {
				fmt.Println(jserr)
			}
			fmt.Println(userMessage)
			rawMessage.ChannelId = channel.Id
			rawMessage.UserId = int(userId)
			rawMessage.Text = userMessage.Text

			err := db.Create(&rawMessage).Error
			if err != nil {
				log.Info("create error:", err)
			}

			dberr := repository.AddUserChannels(userId, channel.Id)

			if dberr != nil {
				log.Fatal(dberr)
			}
			wsResponse.Message = userMessage
			wsResponse.OnlineCount = channelUsers[channelId]
			wsResponseByte, _ := json.Marshal(wsResponse)
			for client := range clients {
				if err = client.WriteMessage(mt, wsResponseByte); err != nil {
					log.Info("write:", err)
					conn.Close()
					break
				}
			}
		}
	}))

	routes.ChannelRoutes(app)
	routes.UserRotues(app)

	log.Fatal(app.Listen(":3001"))

}
