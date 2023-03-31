package main

import (
	"fmt"
	"log"

	"my-finances-api/src/configs"
	"my-finances-api/src/database"
	"my-finances-api/src/handlers"
	"my-finances-api/src/models"

	"github.com/gofiber/fiber/v2"
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

	stockConfigDB := database.DbConfigs{
		Host:     configs.Envs["STOCK_DB_URL"],
		User:     configs.Envs["STOCK_DB_USER"],
		Password: configs.Envs["STOCK_DB_PWD"],
		Name:     configs.Envs["STOCK_DB_NAME"],
		Port:     configs.Envs["STOCK_DB_PORT"],
		SslMode:  configs.Envs["STOCK_DB_SSL"],
		TimeZone: configs.Envs["STOCK_DB_TZ"],
	}

	database.BankDB = bankConfigDB.Open()
	database.Stockdb = stockConfigDB.Open()

	database.BankDB.AutoMigrate(models.Bank{})
	database.BankDB.AutoMigrate(models.Statement{})
	database.BankDB.AutoMigrate(models.User{})
	database.Stockdb.AutoMigrate(models.Bank{})

	app := fiber.New()

	app.Post("/new_user", handlers.CreatNewUser)
	app.Post("/token", handlers.GenerateToken)

	route := app.Group("/secure", handlers.AuthenticateToken)
	{
		route.Get("/stock_db", handlers.GetStocks)
		route.Get("/bank_db", handlers.GetBank)
		route.Get("/stock/:name/price", handlers.GetStockPrice)
		route.Get("/stock/:name/history", handlers.GetStockHistory)
	}

	app.Listen(":3000")
}
