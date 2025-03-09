package handlers

import (
	"Text-Gathering-Service/internal/repository"
	"Text-Gathering-Service/internal/services"
	"Text-Gathering-Service/misc"
	"Text-Gathering-Service/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var clientsMutex = make(map[string]*websocket.Conn)

func _getText(text string, uuid string) (string, bool) {
	var res_text string
	status := false

	services.ClearSpace(&text)
	repo := repository.New()

	if !services.IsLaoText(text) {
		res_text = "ປະໂຫຍກນີ້ບໍ່ແມ່ນພາສາລາວ"
	} else if !services.CheckLaoFormat(text) {
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

func ServeAdminPage(ip string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("admin", fiber.Map{
			"ip_addr": ip,
		})
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

	var res models.ResponseDatas

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Client disconnected: ", uuid)
			delete(clientsMutex, uuid)
			return
		}

		var receivedMsg models.RequestDatas

		if json.Unmarshal(msg, &receivedMsg) != nil {
			res.Content = fmt.Sprintf("Error parsing JSON: %v", err)
			res.Status = false
			continue
		}

		res.Content, res.Status = _getText(receivedMsg.Text, uuid)

		if res.Status {
			audioBytes, err := base64.StdEncoding.DecodeString(receivedMsg.Audio)
			if err != nil {
				res.Content = fmt.Sprintf("Error invalid format: %v", err)
				res.Status = false
				continue
			}

			os.MkdirAll("internal/repository/wait_clips/", os.ModePerm)
			filename := fmt.Sprintf("internal/repository/wait_clips/%s.wav", receivedMsg.Text)
			if os.WriteFile(filename, audioBytes, 0644) != nil {
				res.Content = fmt.Sprintf("Error saving audio: %v", err)
				res.Status = false
				continue
			}
		}

		c.WriteMessage(websocket.TextMessage, misc.Must(json.Marshal(res)))
	}
}

func CheckIncomeDatas(c *websocket.Conn) {
	defer c.Close()

	log.Printf("Admin login\n")

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Admin disconnected")
			return
		}

		if services.AutorizeAdmin(msg) {
			fmt.Println("True: ", msg)
		} else {
			fmt.Println("Wrong: ", msg)
		}
	}
}
