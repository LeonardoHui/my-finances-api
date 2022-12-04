package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type dbConfigs struct {
	host     string
	user     string
	password string
	name     string
	port     string
	sslmode  string
	timeZone string
}

// The struct Properties must be uppercase, so that gorm can access it
type bank struct {
	ID   string
	Name string
}

type stock struct {
	ID   string
	Name string
}

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
	envs, err := godotenv.Read(envFile)

	bankConfigDB := dbConfigs{
		host:     envs["BANK_DB_URL"],
		user:     envs["BANK_DB_USER"],
		password: envs["BANK_DB_PWD"],
		name:     envs["BANK_DB_NAME"],
		port:     envs["BANK_DB_PORT"],
		sslmode:  envs["BANK_DB_SSL"],
		timeZone: envs["BANK_DB_TZ"],
	}

	stockConfigDB := dbConfigs{
		host:     envs["STOCK_DB_URL"],
		user:     envs["STOCK_DB_USER"],
		password: envs["STOCK_DB_PWD"],
		name:     envs["STOCK_DB_NAME"],
		port:     envs["STOCK_DB_PORT"],
		sslmode:  envs["STOCK_DB_SSL"],
		timeZone: envs["STOCK_DB_TZ"],
	}

	// Example: dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	dsn1 := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		bankConfigDB.host, bankConfigDB.user, bankConfigDB.password, bankConfigDB.name, bankConfigDB.port,
	)

	dsn2 := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		stockConfigDB.host, stockConfigDB.user, stockConfigDB.password, stockConfigDB.name, stockConfigDB.port,
	)

	gormConf := gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	bankDB, err := gorm.Open(postgres.Open(dsn1), &gormConf)
	if err != nil {
		log.Fatalln(err)
	}

	stockdb, err := gorm.Open(postgres.Open(dsn2), &gormConf)
	if err != nil {
		log.Fatalln(err)
	}

	app := fiber.New()

	var bank bank
	var stock stock

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
