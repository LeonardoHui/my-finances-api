package handlers

import (
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

func GenerateToken(ctx *fiber.Ctx) error {
	var request TokenRequest
	var user models.User

	if err := ctx.BodyParser(&request); err != nil {
		log.Println("Fail to parse body", err)
		return ERROR_INVALID_PAYLOAD
	}

	if result := database.BankDB.Where("email = ?", request.Email).First(&user); result.Error != nil {
		log.Println("Invalid Email", result.Error)
		return ERROR_INVALID_EMAIL
	}

	if credentialError := user.CheckPassword(request.Password); credentialError != nil {
		log.Println("Invalid Password", credentialError)
		return ERROR_INVALID_PASSWORD
	}

	tokenString, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		log.Println("Error generating JWT", err)
		return ERROR_GENERATING_JWT
	}

	return ctx.JSON(TokenResponse{Token: tokenString})
}

func AuthenticateToken(ctx *fiber.Ctx) error {

	tokenString := ctx.Get("Authorization", "0")
	if tokenString == "0" {
		log.Println("Missing authorization")
		return ERROR_NOT_AUTHORIZED
	}
	user, err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Println("Invalid token:", err)
		return ERROR_NOT_AUTHORIZED
	}
	ctx.Locals("user", user)
	return ctx.Next()
}
