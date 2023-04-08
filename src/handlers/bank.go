package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

type UserBankAccount struct {
	Amount int64  `json:"amount"`
	Bank   string `json:"bank"`
}

func SetBankAccount(ctx *fiber.Ctx) error {
	var user *models.User
	var bankAccount models.BankAccount
	var payload UserBankAccount
	var bank models.Bank

	user = ctx.Locals("user").(*models.User)

	err := json.Unmarshal(ctx.Body(), &payload)
	if err != nil {
		return ERROR_INVALID_PAYLOAD
	}

	err = database.BankDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("name = ?", payload.Bank).First(&bank).Error
		if err == gorm.ErrRecordNotFound {
			log.Println("Bank not found, creating one")
			bank.Name = payload.Bank
			if err = tx.Create(&bank).Error; err != nil {
				log.Println("Fail creating bank, rolling back")
				return err
			}
		}

		result := tx.Model(&models.BankAccount{}).
			Where("user_ID = ? AND bank_ID = ?", user.ID, bank.ID).
			Update("amount", payload.Amount).RowsAffected
		if result == 0 {
			log.Println("Account not found, creating one")
			bankAccount.UserID = user.ID
			bankAccount.Amount = payload.Amount
			bankAccount.BankID = bank.ID
			if err = tx.Create(&bankAccount).Error; err != nil {
				log.Println("Fail creating account, rolling back")
				return err
			}
		}
		return nil
	})

	if err != nil {
		return ERROR_UPDATING_DATA
	}

	return ctx.SendStatus(fiber.StatusOK)
}
