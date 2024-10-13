package main

import (
	//database "chat-app/db"
	"chat-app/routes"
	"chat-app/types"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"strconv"
)

func main() {
	//var upgrader = websocket.Upgrader{
	//	CheckOrigin: func(r *http.Request) bool {
	//		return true
	//	},
	//}
	//var clients = make(map[*websocket.Conn]bool)
	//var broadcast = make(chan types.Message)
	//
	//handleMessages := func() {
	//	for {
	//		msg := <-broadcast
	//		for client := range clients {
	//			err := client.WriteJSON(msg)
	//			if err != nil {
	//				log.Fatal(err)
	//				err := client.Close()
	//				if err != nil {
	//					return
	//				}
	//				delete(clients, client)
	//			}
	//		}
	//	}
	//}
	//
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Welcome to the Chat Room!")
	//})
	//
	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	conn, err := upgrader.Upgrade(w, r, nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer conn.Close()
	//
	//	clients[conn] = true
	//
	//	for {
	//		var msg types.Message
	//		err := conn.ReadJSON(&msg)
	//		if err != nil {
	//			log.Printf("error: %v", err)
	//			delete(clients, conn)
	//			break
	//		}
	//		broadcast <- msg
	//	}
	//})
	//
	//go handleMessages()
	//fmt.Println("Server started on :3000")
	//err := http.ListenAndServe(":3000", nil)
	//if err != nil {
	//	panic(err)
	//}

	app := fiber.New()
	//db := database.GetDB()
	var clients = make(map[*websocket.Conn]types.Channel)
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/:id", websocket.New(func(conn *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		param := conn.Params("id")
		channelId, _ := strconv.ParseUint(param, 10, 32)
		clients[conn] = types.Channel{
			Id: channelId,
		}
		for {
			if mt, msg, err = conn.ReadMessage(); err != nil {
				log.Info("read:", err)
				break
			}
			for client := range clients {
				//conn.WriteMessage(1, msg)
				//log.Info("recv: %s", msg)
				if err = client.WriteMessage(mt, msg); err != nil {
					log.Info("write:", err)
					break
				}
			}
		}

	}))
	routes.ChannelRoutes(app)
	log.Fatal(app.Listen(":3000"))

}
