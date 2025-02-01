package handlers

import (
	"Text-Gathering-Service/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func ServeWebPage() fiber.Handler {
	return static.New("./public")
}

func GetText(c fiber.Ctx) error {
	text := c.FormValue("message")
	status := services.CheckText(text)
	if status {
		return c.JSON(fiber.Map{"message": "Woww what a nice text"})
	} else {
		return c.JSON(fiber.Map{"message": "Nahh I already have that one"})
	}
}
