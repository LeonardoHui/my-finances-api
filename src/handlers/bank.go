package handlers

import (
	"fmt"
	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetBank(ctx *fiber.Ctx) error {
	var bank models.Bank

	database.BankDB.First(&bank)
	return ctx.SendString(fmt.Sprintf("%v", bank))
}

type UserStatements struct {
	Statements []models.Statement   `json:"statements"`
	Balance    []models.BankAccount `json:"balance"`
}

func GetUserStatements(ctx *fiber.Ctx) error {
	var user *models.User
	user = ctx.Locals("user").(*models.User)

	database.BankDB.Preload("Statements").Find(&user)
	database.BankDB.Preload("BankAccounts").Find(&user)
	return ctx.JSON(
		UserStatements{
			Statements: user.Statements,
			Balance:    user.BankAccounts,
		},
	)
}
