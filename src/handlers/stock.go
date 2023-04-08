package handlers

import (
	"fmt"

	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetStocks(ctx *fiber.Ctx) error {
	var stock models.Stock

	database.BankDB.First(&stock)
	return ctx.SendString(fmt.Sprintf("%v", stock))
}

func GetUserInvestments(ctx *fiber.Ctx) error {
	var user *models.User
	user = ctx.Locals("user").(*models.User)

	database.BankDB.Preload("Investments").Find(&user)
	return ctx.JSON(user.Statements)
}
