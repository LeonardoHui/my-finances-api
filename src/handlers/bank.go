package handlers

import (
	"fmt"
	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetBank(c *fiber.Ctx) error {
	var bank models.Bank
	database.BankDB.First(&bank)
	return c.SendString(fmt.Sprintf("%v", bank))
}
