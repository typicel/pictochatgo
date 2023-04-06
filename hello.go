package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	isClosing bool
	mu        sync.Mutex
}

type Room struct {
	Code    string
	Clients map[*websocket.Conn]bool
}

var clients = make(map[*websocket.Conn]*Client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)

// var rooms = make(map[string]*Room)

// func run() {
// 	for {
// 		select {
// 		case connection := <-register:
// 			clients[connection] = &Client{}
// 			log.Println("mmm yummy connection!!!!")

// 		case message := <-broadcast:
// 			log.Println("full of message: ", message)

// 			for connection, c := range clients {
// 				go func(connection *websocket.Conn, c *Client) {
// 					c.mu.Lock()

// 					defer c.mu.Unlock()

// 					if c.isClosing {
// 						return
// 					}

// 					if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
// 						c.isClosing = true
// 						log.Println("BYEEE", err)

// 						connection.WriteMessage(websocket.CloseMessage, []byte{})
// 						connection.Close()
// 						unregister <- connection
// 					}
// 				}(connection, c)
// 			}

// 		case connection := <-unregister:
// 			delete(clients, connection)

// 			log.Println("bye guy!")
// 		}
// 	}
// }

// func handleWebSocket(c *websocket.Conn) {
// 	defer func() {
// 		unregister <- c
// 		c.Close()
// 	}()

// 	register <- c

// 	for {
// 		messageType, message, err := c.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Println("Something went wrong: ", err)
// 			}

// 			return
// 		}

// 		if messageType == websocket.TextMessage {
// 			fmt.Println("Message received")
// 			broadcast <- string(message)
// 		} else {
// 			log.Println("Unknown message type: ", messageType)
// 		}

// 	}
// }

// func createRoom(c *fiber.Ctx) error {
// 	roomCode := c.Params("room")
// 	room := &Room{
// 		Code:      roomCode,
// 		Clients:   make(map[*websocket.Conn]bool),
// 		Broadcast: make(chan string),
// 	}

// 	rooms[roomCode] = room
// 	return c.SendString("room created")
// }

// func joinRoom(c *fiber.Ctx) error {
// 	code := c.Params("code")

// 	room, exists := rooms[code]
// 	if !exists {
// 		return c.Status(404).SendString("No Room with given code")
// 	}

// 	websocket.New(func(conn *websocket.Conn) {
// 		defer func() {
// 			unregister <- conn
// 			conn.Close()
// 		}()

// 		for {
// 			messageType, message, err := conn.ReadMessage()
// 			if err != nil {
// 				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 					log.Println("Something went wrong: ", err)
// 				}
// 				return
// 			}

// 			if messageType == websocket.TextMessage {
// 				room.Broadcast <- string(message)
// 			}

// 		}

// 	})

// 	return nil
// }

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
	app.Static("/", "./public")

	// app.Use(func(c *fiber.Ctx) error {
	// 	if websocket.IsWebSocketUpgrade(c) {
	// 		return c.Next()
	// 	}
	// 	return c.SendStatus(fiber.StatusUpgradeRequired)
	// })

	// go run()
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

	// app.Get("/ws", websocket.New(handleWebSocket))
	// app.Post("/createRoom/:room", createRoom)
	// app.Get("/room/:code", joinRoom)

	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))

}
