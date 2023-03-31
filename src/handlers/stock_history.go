package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"my-finances-api/src/configs"

	"github.com/gofiber/fiber/v2"
)

func GetStockHistory(c *fiber.Ctx) error {

	req := http.Client{}
	stockName := c.Params("name")
	url := fmt.Sprintf(
		"%s/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&outputsize=full&apikey=%s",
		configs.Envs["AP_API_URL"],
		stockName,
		configs.Envs["AP_API_KEY"],
	)
	resp, err := req.Get(url)
	if err != nil {
		return c.SendString("Error: Try again later")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	return c.SendString(string(respBody))
}
