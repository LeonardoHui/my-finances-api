package handlers

import (
	"log"
	"my-finances-api/src/auth"
	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

func CreatNewUser(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		log.Println("Error parsing body", err)
		return ERROR_INVALID_PAYLOAD
	}
	if err := user.HashPassword(user.Password); err != nil {
		log.Println("Error hashing password", err)
		return ERROR_CREATING_USER
	}
	record := database.BankDB.Create(&user)
	if record.Error != nil {
		log.Println("Error saving DB", record.Error)
		return ERROR_CREATING_USER
	}

	if err := ctx.JSON(user); err != nil {
		log.Println("Error returning body", err)
		return ERROR_CREATING_USER
	}
	return nil
}

func CreatNewUserAndLogin(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		log.Println("Error parsing body", err)
		return ERROR_INVALID_PAYLOAD
	}

	if err := CreatNewUser(ctx); err != nil {
		log.Println("Error Creating New User", err)
		return ERROR_CREATING_USER
	}

	tokenString, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		log.Println("Error generating JWT", err)
		return ERROR_GENERATING_JWT
	}

	return ctx.Status(fiber.StatusCreated).JSON(TokenResponse{Token: tokenString})
}
