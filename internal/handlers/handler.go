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
	"github.com/google/uuid"
)

var clientsMutex = make(map[string]*websocket.Conn)

func _getText(text string, uuid string) (string, bool) {
	var res_text string
	status := false

	services.ClearSpace(&text)
	isLao := services.IsLaoText(text)
	isFormatError := services.CheckLaoFormat(text)
	repo := repository.New()

	if !isLao {
		res_text = "ປະໂຫຍກນີ້ບໍ່ແມ່ນພາສາລາວ"
	} else if !isFormatError {
		res_text = "ຮູບແບບຂອງປະໂຫຍກບໍ່ຖຶກຕ້ອງ"
	} else if !repo.StoreIntoDB(text) {
		res_text = "ປະໂຫຍກນີ້ມີໃນລະບົບແລ້ວກະລຸນາປ້ອນໃຫມ່"
	} else {
		res_text = "ພວກເຮົາໄດ້ບັນທຶກປະໂຫຍກທີ່ຖຶກປ້ອນໄວ້ໃນລະບົບແລ້ວ"
		status = true
	}

	log.Printf("UUID:[%s] input: %s | Status: %v\n", uuid, text, status)

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

	uuid := uuid.New().String()
	clientsMutex[uuid] = c

	log.Printf("Client UUID:[%s] connected\n", uuid)

	repo := repository.New()
	c.WriteMessage(websocket.TextMessage, misc.Must(json.Marshal(repo.GetAllCategoryDatas())))

	log.Printf("UUID:[%s]: -- Update categories --\n", uuid)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Client disconnected: ", uuid)
			delete(clientsMutex, uuid)
			return
		}

		t, s := _getText(string(msg), uuid)
		res := models.ResponseDatas{
			Content: t,
			Status:  s,
		}

		c.WriteMessage(websocket.TextMessage, misc.Must(json.Marshal(res)))
	}

}
