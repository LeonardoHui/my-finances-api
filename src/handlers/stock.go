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

type UserInvestments struct {
	Investments       []models.InvestmentEvent `json:"investments"`
	StockDistribution []models.Investment      `json:"stock_distribution"`
	DividendYield     []models.InvestmentEvent `json:"dividend_yied"`
	DidivendPaid      []models.InvestmentEvent `json:"dividend_paid"`
}

func GetUserInvestments(ctx *fiber.Ctx) error {
	var user *models.User
	user = ctx.Locals("user").(*models.User)

	database.BankDB.Preload("Investments").Find(&user)
	database.BankDB.Preload("InvestmentEvents").Find(&user)
	return ctx.JSON(
		UserInvestments{
			Investments:       user.InvestmentEvents,
			StockDistribution: user.Investments,
			DividendYield:     []models.InvestmentEvent{},
			DidivendPaid:      []models.InvestmentEvent{},
		},
	)
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
