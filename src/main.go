package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"my-finances-api/src/auth"
	"my-finances-api/src/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
	"gorm.io/gorm"
)

var (
	BankDB  *gorm.DB
	Stockdb *gorm.DB
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateToken(context *fiber.Ctx) error {
	var request TokenRequest
	var user database.User

	if err := context.BodyParser(&request); err != nil {
		log.Println("Fail to parse boday", err)
		return err
	}

	if result := BankDB.Where("email = ?", request.Email).First(&user); result.Error != nil {
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

	return context.JSON(tokenString)
}

func Auth(context *fiber.Ctx) error {

	tokenString := context.Get("Authorization", "0")
	if tokenString == "0" {
		log.Println("Não autorizado - Missing authorization")
		return errors.New("Não autorizado - Missing authorization")
	}
	err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Println("Invalid token:", err)
		return err
	}
	return context.Next()
}

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

	envFile := os.Args[1]
	envs, _ := godotenv.Read(envFile)

	bankConfigDB := database.DbConfigs{
		Host:     envs["BANK_DB_URL"],
		User:     envs["BANK_DB_USER"],
		Password: envs["BANK_DB_PWD"],
		Name:     envs["BANK_DB_NAME"],
		Port:     envs["BANK_DB_PORT"],
		SslMode:  envs["BANK_DB_SSL"],
		TimeZone: envs["BANK_DB_TZ"],
	}

	stockConfigDB := database.DbConfigs{
		Host:     envs["STOCK_DB_URL"],
		User:     envs["STOCK_DB_USER"],
		Password: envs["STOCK_DB_PWD"],
		Name:     envs["STOCK_DB_NAME"],
		Port:     envs["STOCK_DB_PORT"],
		SslMode:  envs["STOCK_DB_SSL"],
		TimeZone: envs["STOCK_DB_TZ"],
	}

	BankDB = bankConfigDB.Open()
	Stockdb = stockConfigDB.Open()

	BankDB.AutoMigrate(database.Bank{})
	BankDB.AutoMigrate(database.Statement{})
	BankDB.AutoMigrate(database.User{})
	Stockdb.AutoMigrate(database.Bank{})

	app := fiber.New()

	var bank database.Bank
	var stock database.Stock

	app.Post("/new_user", func(c *fiber.Ctx) error {
		var user database.User
		if err := c.BodyParser(&user); err != nil {
			log.Println("Error parsing body", err)
			return err
		}
		if err := user.HashPassword(user.Password); err != nil {
			log.Println("Error hashing password", err)
			return err
		}
		record := BankDB.Create(&user)
		if record.Error != nil {
			log.Println("Error saving DB", record.Error)
			return record.Error
		}
		if err := c.JSON(user); err != nil {
			log.Println("Error returning body", err)
			return err
		}
		return nil
	})

	app.Post("/token", GenerateToken)

	route := app.Group("/secure", Auth)
	route.Get("/stock_db", func(c *fiber.Ctx) error {
		Stockdb.First(&stock)
		return c.SendString(fmt.Sprintf("%v", stock))
	})

	route.Get("/bank_db", func(c *fiber.Ctx) error {
		BankDB.First(&bank)
		return c.SendString(fmt.Sprintf("%v", bank))
	})

	route.Get("/stock/:name/price", func(c *fiber.Ctx) error {
		r := mockrequest()
		req := http.Client{Transport: r}
		defer r.Stop() // Make sure recorder is stopped once done with it
		//Example of ThirdParty API request
		//https://api.hgbrasil.com/finance/stock_price?key=c34b53c1&symbol=mxrf11
		stockName := c.Params("name")
		url := fmt.Sprintf("https://api.hgbrasil.com/finance/stock_price?symbol=%s&key=%s", stockName, envs["HG_API_KEY"])
		resp, err := req.Get(url)
		if err != nil {
			return c.SendString("Error: Try again later")
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		return c.SendString(string(respBody))
	})

	route.Get("/stock/:name/history", func(c *fiber.Ctx) error {
		r := mockrequest()
		req := http.Client{Transport: r}
		defer r.Stop() // Make sure recorder is stopped once done with it
		//Example of ThirdParty API request
		//https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=IBM&outputsize=full&apikey=FIIX3DBF99Y439X5
		stockName := c.Params("name")
		url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&outputsize=full&apikey=%s", stockName, envs["AP_API_KEY"])
		resp, err := req.Get(url)
		if err != nil {
			return c.SendString("Error: Try again later")
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		return c.SendString(string(respBody))
	})

	app.Listen(":3000")
}
