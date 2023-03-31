package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"my-finances-api/src/configs"

	"github.com/gofiber/fiber/v2"
)

func GetStockPrice(c *fiber.Ctx) error {

	req := http.Client{}
	stockName := c.Params("name")
	url := fmt.Sprintf(
		"%s/finance/stock_price?symbol=%s&key=%s",
		configs.Envs["HG_API_URL"],
		stockName,
		configs.Envs["HG_API_KEY"],
	)
	resp, err := req.Get(url)
	if err != nil {
		return c.SendString("Error: Try again later")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	return c.SendString(string(respBody))
}
