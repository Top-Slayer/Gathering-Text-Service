package handlers

import (
	"Text-Gathering-Service/internal/repository"
	"Text-Gathering-Service/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func ServeWebPage() fiber.Handler {
	return static.New("./public")
}

func GetText(c fiber.Ctx) error {
	var res_text string
	color_status := false
	text := c.FormValue("message")

	isLao := services.IsLaoText(text)
	repo := repository.New()

	if !isLao {
		res_text = "That's not Lao langages"
	} else if !repo.StoreIntoDB(text) {
		res_text = "Datas already have"
	} else {
		res_text = "Thank's you for helping"
		color_status = true
	}

	return c.JSON(fiber.Map{"message": res_text, "status_color": color_status})
}
