package handlers

import (
	"fmt"

	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetStocks(c *fiber.Ctx) error {
	var stock models.Stock
	database.Stockdb.First(&stock)
	return c.SendString(fmt.Sprintf("%v", stock))
}
