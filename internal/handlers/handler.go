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
		res_text = "That's not Lao langages"
	} else if !repo.StoreIntoDB(text) {
		res_text = "Datas already have"
	} else {
		res_text = "Thank's you for helping"
		status = true
	}

	log.Printf("Incoming message from client: %s | Status: %v\n", text, status)

	return res_text, status
}

func UpgradeWebsocketProtocol(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		log.Println("Client connected")
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func ConnectWebsocket(c *websocket.Conn) {
	defer c.Close()

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
