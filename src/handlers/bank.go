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

type UserStatements struct {
	Statements []models.Statement   `json:"statements"`
	Balance    []models.BankAccount `json:"balance"`
}

func GetUserStatements(c *fiber.Ctx) error {
	var user *models.User
	user = c.Locals("user").(*models.User)

	database.BankDB.Preload("Statements").Find(&user)
	database.BankDB.Preload("BankAccounts").Find(&user)
	return c.JSON(
		UserStatements{
			Statements: user.Statements,
			Balance:    user.BankAccounts,
		},
	)
}
