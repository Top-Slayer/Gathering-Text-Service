package main

import (
	"Text-Gathering-Service/internal/handlers"
	"Text-Gathering-Service/misc"
	"io"
	"log"
	"os"

	"github.com/gofiber/template/html/v2"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func _createLogFIle() {
	for {
		file := misc.Must(os.OpenFile("details.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666))
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	}
}

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

	// app.Use("/send-text", handlers.UpgradeWebsocketProtocol)

	app.Get("/send-text", websocket.New(handlers.GetDatasFromClient))

	go _createLogFIle()
	log.Fatal(app.Listen(":3000"))
}
