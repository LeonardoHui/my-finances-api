package handlers

import (
	"errors"
	"log"
	"my-finances-api/src/auth"
	"my-finances-api/src/database"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func GenerateToken(context *fiber.Ctx) error {
	var request TokenRequest
	var user models.User

	if err := context.BodyParser(&request); err != nil {
		log.Println("Fail to parse boday", err)
		return err
	}

	if result := database.BankDB.Where("email = ?", request.Email).First(&user); result.Error != nil {
		log.Println("Invalid Email", result.Error)
		return result.Error
	}

	if credentialError := user.CheckPassword(request.Password); credentialError != nil {
		log.Println("Invalid Password", credentialError)
		return credentialError
	}

	tokenString, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		return err
	}

	return context.JSON(TokenResponse{Token: tokenString})
}

func AuthenticateToken(context *fiber.Ctx) error {

	tokenString := context.Get("Authorization", "0")
	if tokenString == "0" {
		log.Println("Não autorizado - Missing authorization")
		return errors.New("Não autorizado - Missing authorization")
	}
	user, err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Println("Invalid token:", err)
		return err
	}
	context.Locals("user", user)
	return context.Next()
}
