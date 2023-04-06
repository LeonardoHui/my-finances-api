package handlers

import (
	"log"
	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func CreatNewUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Println("Error parsing body", err)
		return err
	}
	if err := user.HashPassword(user.Password); err != nil {
		log.Println("Error hashing password", err)
		return err
	}
	record := database.BankDB.Create(&user)
	if record.Error != nil {
		log.Println("Error saving DB", record.Error)
		return record.Error
	}

	if err := c.JSON(user); err != nil {
		log.Println("Error returning body", err)
		return err
	}
	return nil
}

// For teste only
func InternalCreateNewUser() {
	user := models.User{
		Name:     "tester one",
		Username: "tester one",
		Email:    "tester_one@mail.com",
		Password: "1111",
	}
	if err := user.HashPassword(user.Password); err != nil {
		log.Println("Error hashing password", err)

	}
	record := database.BankDB.Create(&user)
	if record.Error != nil {
		log.Println("Error saving DB", record.Error)
	}
}
