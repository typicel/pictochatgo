package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type client struct {
	isClosing bool
	mu        sync.Mutex
}

var clients = make(map[*websocket.Conn]*client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)

func run() {
	for {
		select {
		case connection := <-register:
			clients[connection] = &client{}
			log.Println("mmm yummy connection!!!!")

		case message := <-broadcast:
			log.Println("full of message: ", message)

			for connection, c := range clients {
				go func(connection *websocket.Conn, c *client) {
					c.mu.Lock()

					defer c.mu.Unlock()

					if c.isClosing {
						return
					}

					if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
						c.isClosing = true
						log.Println("BYEEE", err)

						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
						unregister <- connection
					}
				}(connection, c)
			}

		case connection := <-unregister:
			delete(clients, connection)

			log.Println("bye guy!")
		}
	}
}

func main() {
	app := fiber.New()

	app.Static("/", "./public")

	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	go run()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c
			c.Close()
		}()

		register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("guh!: ", err)
				}

				return
			}

			if messageType == websocket.TextMessage {
				fmt.Println("wow nice message!")
				broadcast <- string(message)
			} else {
				log.Println("??????????????????????", messageType)
			}

		}
	}))

	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))

}
