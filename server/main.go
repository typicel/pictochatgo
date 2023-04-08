package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Room struct {
	Code    string
	Clients map[*websocket.Conn]bool
}

func (r *Room) addClient(c *websocket.Conn) {
	r.Clients[c] = true
}

func (r *Room) removeClient(c *websocket.Conn) {
	delete(r.Clients, c)
}

func (r *Room) broadcast(message []byte) {
	for client := range r.Clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("error broadcasting: ", err)
			r.removeClient(client)
		}
	}
}

var rooms []*Room

func main() {
	app := fiber.New()

	// post route to create a new room
	app.Post("/create/:room", func(c *fiber.Ctx) error {
		roomCode := c.Params("room")
		fmt.Println("Creating room: ", roomCode)

		room := &Room{
			Code:    roomCode,
			Clients: make(map[*websocket.Conn]bool),
		}

		rooms = append(rooms, room)
		return c.SendStatus(200)
	})

	app.Get("/ws/:room", websocket.New(func(c *websocket.Conn) {
		roomName := c.Params("room")
		fmt.Println("Joining Room: ", roomName)

		var currentRoom *Room
		for _, r := range rooms {
			if r.Code == roomName {
				fmt.Println("Room found: ", r.Code)
				currentRoom = r
				break
			}
		}

		if currentRoom == nil {
			return
		}

		currentRoom.addClient(c)

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("error reading message:", err)
				currentRoom.removeClient(c)
				break
			}
			if messageType == websocket.TextMessage {
				currentRoom.broadcast(message)
			}
		}
	}))

	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))

}
