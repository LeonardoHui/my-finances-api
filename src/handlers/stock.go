package handlers

import (
	"fmt"

	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetStocks(c *fiber.Ctx) error {
	var stock models.Stock

	database.BankDB.First(&stock)
	return c.SendString(fmt.Sprintf("%v", stock))
}

func GetUserInvestments(c *fiber.Ctx) error {
	var user *models.User
	user = c.Locals("user").(*models.User)

	database.BankDB.Preload("Investments").Find(&user)
	return c.JSON(user.Statements)
}
