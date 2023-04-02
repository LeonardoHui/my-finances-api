package main

import (
	"fmt"
	"log"

	"my-finances-api/src/configs"
	"my-finances-api/src/database"
	"my-finances-api/src/handlers"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func mockrequest() *recorder.Recorder {
	r, err := recorder.New("fixtures/integration")
	if err != nil {
		log.Fatal(err)
	}

	if r.Mode() != recorder.ModeRecordOnce {
		log.Fatal("Recorder should be in ModeRecordOnce")
	}
	return r
}

func main() {
	fmt.Println("STARTING THE PROGRAM")

	bankConfigDB := database.DbConfigs{
		Host:     configs.Envs["BANK_DB_URL"],
		User:     configs.Envs["BANK_DB_USER"],
		Password: configs.Envs["BANK_DB_PWD"],
		Name:     configs.Envs["BANK_DB_NAME"],
		Port:     configs.Envs["BANK_DB_PORT"],
		SslMode:  configs.Envs["BANK_DB_SSL"],
		TimeZone: configs.Envs["BANK_DB_TZ"],
	}

	database.BankDB = bankConfigDB.Open()

	database.BankDB.AutoMigrate(models.Bank{})
	database.BankDB.AutoMigrate(models.Statement{})
	database.BankDB.AutoMigrate(models.User{})
	database.BankDB.AutoMigrate(models.Stock{})

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
	}))

	app.Post("/register", handlers.CreatNewUser)
	app.Post("/login", handlers.GenerateToken)

	// Authenticated endpoints
	app.Get("/investments", handlers.AuthenticateToken, handlers.GetStocks)
	app.Get("/statements", handlers.AuthenticateToken, handlers.GetBank)

	app.Listen(":8000")
}
