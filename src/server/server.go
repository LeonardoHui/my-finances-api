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

	app.Post("/register", handlers.RequestLogger, handlers.CreatNewUserAndLogin)
	app.Post("/login", handlers.RequestLogger, handlers.GenerateToken)
	app.Get("/investments", handlers.RequestLogger, handlers.GetUserInvestmentsEvolution)
	app.Get("/simulation", handlers.RequestLogger, handlers.GetInvestmentsEvolutionSimulation)
	// Testing
	app.Get("/distint", handlers.RequestLogger, handlers.GetStocks)
	app.Get("/market", handlers.RequestLogger, handlers.GetMarketValueTestHandler)
	app.Get("/timeline", handlers.RequestLogger, handlers.GetTimeLine)

	// Authenticated endpoints
	app.Get("/statements", handlers.RequestLogger, handlers.AuthenticateToken, handlers.GetUserStatements)
	app.Post("/bank/account", handlers.RequestLogger, handlers.AuthenticateToken, handlers.SetBankAccount)
	app.Post("/bank/events", handlers.RequestLogger, handlers.AuthenticateToken, handlers.SetBankEvent)
	app.Post("/stock", handlers.RequestLogger, handlers.AuthenticateToken, handlers.SetInvestment)
	app.Post("/stock/events", handlers.RequestLogger, handlers.AuthenticateToken, handlers.SetInvestmentEvent)

	app.Listen(Port)
}
