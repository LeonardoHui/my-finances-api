package server

import (
	"my-finances-api/src/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Run(Port string) {

	config := fiber.Config{ErrorHandler: handlers.ResponseWhenError}
	app := fiber.New(config)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
	}))

	app.Post("/register", handlers.CreatNewUserAndLogin)
	app.Post("/login", handlers.GenerateToken)
	app.Get("/investments", handlers.GetUserInvestmentsEvolution)
	app.Get("/simulation", handlers.GetInvestmentsEvolutionSimulation)

	// Authenticated endpoints
	app.Get("/statements", handlers.AuthenticateToken, handlers.GetUserStatements)
	app.Post("/bank/account", handlers.AuthenticateToken, handlers.SetBankAccount)
	app.Post("/bank/events", handlers.AuthenticateToken, handlers.SetBankEvent)
	app.Post("/stock", handlers.AuthenticateToken, handlers.SetInvestment)
	app.Post("/stock/events", handlers.AuthenticateToken, handlers.SetInvestmentEvent)

	app.Listen(Port)
}
