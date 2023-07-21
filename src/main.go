package main

import (
	"fmt"

	"my-finances-api/src/configs"
	"my-finances-api/src/database"
	"my-finances-api/src/server"
	"my-finances-api/src/utils"
)

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
	database.Migrate()

	// For testing and development
	if configs.Envs["ENV"] != "PROD" {
		utils.InternalCreateNewUser()
		utils.InternalLoadTables("./sql/")
	}

	server.Run(":8000")
}
