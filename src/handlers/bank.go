package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"my-finances-api/src/database"
	"my-finances-api/src/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetBank(ctx *fiber.Ctx) error {
	var bank models.Bank

	database.BankDB.First(&bank)
	return ctx.SendString(fmt.Sprintf("%v", bank))
}

type GenericMonetaryItem struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Amount      uint      `json:"amount"`
	Date        time.Time `json:"date"`
}

type ApiUserStatements struct {
	Statements []GenericMonetaryItem `json:"statements"`
	Balance    []GenericMonetaryItem `json:"balance"`
}

func GetUserStatements(ctx *fiber.Ctx) error {

	var user *models.User
	user = ctx.Locals("user").(*models.User)

	database.BankDB.Preload("Statements").Find(&user)
	database.BankDB.Preload("BankAccounts").Find(&user)

	//Convert from DB to API response
	var (
		balances        = []GenericMonetaryItem{}
		statements      = []GenericMonetaryItem{}
		transactionType string
		bankNameMap     = make(map[uint]string)
	)
	bankNameMap[1] = "BANK ONE"
	bankNameMap[2] = "BANK TWO"
	bankNameMap[3] = "BANK THREE"

	for _, v := range user.Statements {
		if v.Event == "BUY" {
			transactionType = "DEBIT"
		} else {
			transactionType = "CREDIT"
		}
		statements = append(statements, GenericMonetaryItem{
			ID:          v.ID,
			Description: transactionType,
			Amount:      uint(v.Amount),
			Date:        v.CreatedAt,
		})
	}

	for _, v := range user.BankAccounts {
		balances = append(balances, GenericMonetaryItem{
			ID:          v.ID,
			Description: bankNameMap[v.BankID],
			Amount:      uint(v.Amount),
			Date:        v.UpdatedAt,
		})
	}

	return ctx.JSON(
		ApiUserStatements{
			Statements: statements,
			Balance:    balances,
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

func SetBankEvent(ctx *fiber.Ctx) error {
	var user *models.User
	var payload models.Statement

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
