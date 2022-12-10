package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"my-finances-api/src/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func mockrequest() *recorder.Recorder {
	r, err := recorder.New("fixtures/integration")
	if err != nil {
		log.Fatal(err)
	}
	// defer r.Stop() // Make sure recorder is stopped once done with it

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

	bankDB := bankConfigDB.Open()
	stockdb := stockConfigDB.Open()

	app := fiber.New()

	var bank database.Bank
	var stock database.Stock

	app.Get("/stock_db", func(c *fiber.Ctx) error {
		stockdb.First(&stock)
		return c.SendString(fmt.Sprintf("%v", stock))
	})

	app.Get("/bank_db", func(c *fiber.Ctx) error {
		bankDB.First(&bank)
		return c.SendString(fmt.Sprintf("%v", bank))
	})

	app.Get("/stock/:name/price", func(c *fiber.Ctx) error {
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

	app.Get("/stock/:name/history", func(c *fiber.Ctx) error {
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
