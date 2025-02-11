package handlers

import (
	"Text-Gathering-Service/internal/repository"
	"Text-Gathering-Service/internal/services"
	"Text-Gathering-Service/misc"
	"Text-Gathering-Service/models"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func _getText(text string) (string, bool) {
	var res_text string
	status := false

	isLao := services.IsLaoText(text)
	repo := repository.New()

	if !isLao {
		res_text = "ປະໂຫຍກນີ້ບໍ່ແມ່ນພາສາລາວ"
	} else if !repo.StoreIntoDB(text) {
		res_text = "ປະໂຫຍກນີ້ມີໃນລະບົບແລ້ວກະລຸນາປ້ອນໃຫມ່"
	} else {
		res_text = "ພວກເຮົາໄດ້ບັນທຶກປະໂຫຍກທີ່ຖຶກປ້ອນໄວ້ໃນລະບົບແລ້ວ"
		status = true
	}

	log.Printf("Incoming message from client: %s | Status: %v\n", text, status)

	return res_text, status
}

func ServeWebpage(ip string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{"ip_addr": ip})
	}
}

func UpgradeWebsocketProtocol(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func GetDatasFromClient(c *websocket.Conn) {
	defer c.Close()
	log.Printf("Client IP: %s connected\n", c.RemoteAddr().String())

	repo := repository.New()
	c.WriteMessage(websocket.TextMessage, misc.Must(json.Marshal(repo.GetAllCategoryDatas())))

	log.Println("-- Update categories --")

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Client disconnected: ", err)
			return
		}

		t, s := _getText(string(msg))
		res := models.ResponseDatas{
			Content: t,
			Status:  s,
		}

		c.WriteMessage(websocket.TextMessage, misc.Must(json.Marshal(res)))
	}

}
