package server

import (
	"log"
	"my-finances-api/src/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RequestLogger(ctx *fiber.Ctx) error {
	log.Println(ctx.Path())
	return ctx.Next()
}

func Run(Port string) {

	config := fiber.Config{ErrorHandler: handlers.ResponseWhenError}
	app := fiber.New(config)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
	}))

	app.Post("/register", RequestLogger, handlers.CreatNewUserAndLogin)
	app.Post("/login", RequestLogger, handlers.GenerateToken)

	// Authenticated endpoints
	app.Get("/statements", RequestLogger, handlers.AuthenticateToken, handlers.GetUserStatements)
	app.Get("/investments", RequestLogger, handlers.AuthenticateToken, handlers.GetUserInvestments)
	app.Get("/evolution", RequestLogger, handlers.AuthenticateToken, handlers.GetUserInvestmentsEvolution)
	app.Get("/simulation", RequestLogger, handlers.AuthenticateToken, handlers.GetInvestmentsEvolutionSimulation)

	app.Post("/bank/account", RequestLogger, handlers.AuthenticateToken, handlers.SetBankAccount)
	app.Post("/bank/events", RequestLogger, handlers.AuthenticateToken, handlers.SetBankEvent)
	app.Post("/stock", RequestLogger, handlers.AuthenticateToken, handlers.SetInvestment)
	app.Post("/stock/events", RequestLogger, handlers.AuthenticateToken, handlers.SetInvestmentEvent)

	app.Listen(Port)
}
