package handlers

import (
	"fmt"

	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetStocks(c *fiber.Ctx) error {
	var stock models.Stock

	user := (c.Locals("user"))
	fmt.Printf("LOCALS %v", user)

	database.BankDB.First(&stock)
	return c.SendString(fmt.Sprintf("%v", stock))
}
