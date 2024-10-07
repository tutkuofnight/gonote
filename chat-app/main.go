package main

import (
	"chat-app/types"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan types.Message)
	handleMessages := func() {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Fatal(err)
				err := client.Close()
				if err != nil {
					return
				}
				delete(clients, client)
			}
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the Chat Room!")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		clients[conn] = true

		for {
			var msg types.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
			}
			broadcast <- msg
		}
	})

	go handleMessages()
	fmt.Println("Server started on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}

}
