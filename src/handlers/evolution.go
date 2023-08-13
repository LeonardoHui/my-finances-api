package handlers

import (
	"fmt"
	"my-finances-api/src/database"
	"my-finances-api/src/models"
	"my-finances-api/src/utils"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

//Get stocks index/names owned by the user
//For each stock calculate its market valued through time
//Sum all stocks through time

func GetStocks(ctx *fiber.Ctx) error {
	// var user *models.User
	// user = ctx.Locals("user").(*models.User)
	// user_id = user.ID
	var user_id = 1
	var result []string

	database.BankDB.Model(models.StockEvent{}).
		Joins("JOIN stocks ON stocks.id = stock_events.stock_id").
		Where("user_id = ?", user_id).
		Distinct("stock_id").
		Select("name").
		Scan(&result)

	return ctx.SendString(fmt.Sprintf("%v", result))
}

func getMarketValue(stockID uint) map[time.Time]uint {
	// access stock_history filter all history for that stock
	// access user stock events. get all usr history of that stock
	// Map user stock history by quantity and date
	// create a time line of quantity multiplied by the marked history

	var stockHistory []models.StockHistory
	var stockEvents []models.StockEvent
	var user_id = 1

	stocksQtdHistory := make(map[time.Time]uint)
	stocksMarketValue := make(map[time.Time]uint)

	database.BankDB.Order("created_at").Find(&stockEvents, "user_id = ? AND stock_id = ?", user_id, stockID)

	startingDate := stockEvents[0].CreatedAt.UTC().AddDate(0, 0, -1)
	database.BankDB.Order("created_at").Find(&stockHistory, "created_at > ? AND stock_id = ?", startingDate, stockID)

	var sum uint
	for _, v := range stockEvents {
		sum = sum + v.Quantity
		year, month, _ := v.CreatedAt.UTC().Date()
		finalDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
		stocksQtdHistory[finalDate] = sum
	}

	// Check if there are stock events if not use the previous quantity
	var qty = stocksQtdHistory[stockHistory[0].CreatedAt]
	stocksMarketValue[stockHistory[0].CreatedAt] = stocksQtdHistory[stockHistory[0].CreatedAt] * stockHistory[0].Price
	for _, v := range stockHistory {
		if _, ok := stocksQtdHistory[v.CreatedAt]; ok {
			qty = stocksQtdHistory[v.CreatedAt]
		}
		stocksMarketValue[v.CreatedAt] = qty * v.Price
	}

	return stocksMarketValue
}

func GetMarketValueTestHandler(ctx *fiber.Ctx) error {

	type temp struct {
		date  time.Time
		value uint
	}

	result := getMarketValue(1)
	var history []temp
	for i, v := range result {
		history = append(history, temp{value: v, date: i})
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].date.Before(history[j].date)
	})

	return ctx.SendString(fmt.Sprintf("%+v", history))
}

func getOrverallMarketValue() []utils.Timeline {

	var timeline []utils.Timeline

	var user_id = 1
	var firstEvent models.StockEvent
	var stocksList []uint
	var stocksHistoryArray []map[time.Time]uint

	database.BankDB.Model(models.StockEvent{}).
		Joins("JOIN stocks ON stocks.id = stock_events.stock_id").
		Where("user_id = ?", user_id).
		Distinct("stock_id").
		Select("stock_id").
		Scan(&stocksList)

	for _, stock := range stocksList {
		stocksHistoryArray = append(stocksHistoryArray, getMarketValue(stock))
	}

	//Create an array of months starting fron the first users event
	database.BankDB.Order("Created_at").First(&firstEvent, "user_id = ?", user_id)
	startDate := firstEvent.CreatedAt
	timeline = utils.MonthsArray(startDate)

	for i, time := range timeline {
		for _, stock := range stocksHistoryArray {
			timeline[i].Value = timeline[i].Value + stock[time.Date]
		}
	}
	return timeline
}

func GetTimeLine(ctx *fiber.Ctx) error {

	result := getOrverallMarketValue()

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return ctx.JSON(result)
}
