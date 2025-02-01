package main

import (
	"Text-Gathering-Service/internal/handlers"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", handlers.ServeWebPage())

	app.Post("/send-text", handlers.GetText)

	log.Fatal(app.Listen(":3000"))
}
