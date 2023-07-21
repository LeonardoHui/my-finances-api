package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetStocks(ctx *fiber.Ctx) error {
	var stock models.Stock

	database.BankDB.First(&stock)
	return ctx.SendString(fmt.Sprintf("%v", stock))
}

func SetInvestmentEvent(ctx *fiber.Ctx) error {
	var user *models.User
	var payload models.InvestmentEvent

	err := json.Unmarshal(ctx.Body(), &payload)
	if err != nil {
		return ERROR_INVALID_PAYLOAD
	}

	user = ctx.Locals("user").(*models.User)
	payload.UserID = user.ID

	if err := database.BankDB.Create(&payload).Error; err != nil {
		log.Println("Fail creating statement event")
		return ERROR_UPDATING_DATA
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func SetInvestment(ctx *fiber.Ctx) error {
	var user *models.User
	var payload models.Investment

	err := json.Unmarshal(ctx.Body(), &payload)
	if err != nil {
		return ERROR_INVALID_PAYLOAD
	}

	user = ctx.Locals("user").(*models.User)
	payload.UserID = user.ID

	if err := database.BankDB.Create(&payload).Error; err != nil {
		log.Println("Fail creating statement event")
		return ERROR_UPDATING_DATA
	}

	return ctx.SendStatus(fiber.StatusOK)
}
