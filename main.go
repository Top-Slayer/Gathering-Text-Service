package main

import (
	"Text-Gathering-Service/internal/handlers"
	"log"
	"os"

	"github.com/gofiber/template/html/v2"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func _setDomain() string {
	args := os.Args
	var ip string
	if len(args) <= 1 {
		ip = "http://localhost:3000/send-text"
	} else {
		ip = "https://" + args[1] + "/ws"
	}

	log.Println("Server Start: \033[34m" + ip + "\033[0m")

	return ip
}

func main() {
	engine := html.New("./public", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())

	// app.Get("/", handlers.ServeWebpage(_setDomain()))

	app.Use("/send-text", handlers.UpgradeWebsocketProtocol)

	app.Get("/send-text", websocket.New(handlers.GetDatasFromClient))

	log.Fatal(app.Listen(":3000"))
}
