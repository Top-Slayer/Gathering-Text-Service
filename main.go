package main

import (
	"Text-Gathering-Service/internal/handlers"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()
	app.Use(cors.New())

	app.Static("/", "./public")

	app.Use("/ws", handlers.UpgradeWebsocketProtocol)
	app.Get("/ws", websocket.New(handlers.ConnectWebsocket))

	log.Fatal(app.Listen(":3000"))
}
