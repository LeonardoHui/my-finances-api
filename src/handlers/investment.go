package handlers

import (
	"my-finances-api/src/database"
	"my-finances-api/src/models"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetUserInvestments(ctx *fiber.Ctx) error {
	var (
		user        *models.User
		stockEvents []models.StockEvent

		investments    []models.GenericMonetaryItem
		dividendPaid   []models.GenericLabelValue
		dividendYield  []models.GenericLabelValue
		investmentName = make(map[uint]string)
	)
	investmentName[1] = "STOCK A"
	investmentName[2] = "STOCK B"
	investmentName[3] = "STOCK C"
	user = ctx.Locals("user").(*models.User)

	database.BankDB.Order("created_at").Find(&stockEvents, "user_id = ?", user.ID)

	for _, v := range stockEvents {
		// TODO - Sum each stock
		investments = append(investments, models.GenericMonetaryItem{
			ID:          v.ID,
			Description: investmentName[v.StockID],
			Amount:      v.TotalPrice,
			Date:        v.CreatedAt,
		})
	}

	for _, v := range stockEvents {
		if v.Event == "DIVIDEND" {
			dividendPaid = append(dividendPaid, models.GenericLabelValue{
				Label: investmentName[v.StockID],
				Value: v.TotalPrice,
			})
			// TODO - Calculate Yield
			dividendYield = append(dividendYield, models.GenericLabelValue{
				Label: investmentName[v.StockID],
				Value: v.TotalPrice,
			})
		}
	}

	return ctx.JSON(
		models.UserInvestments{
			Investments:   investments,
			DividendYield: dividendYield,
			DividendPaid:  dividendPaid,
		},
	)
}

func GetUserInvestmentsEvolution(ctx *fiber.Ctx) error {
	// Get stocks and bonds from users
	// Pile up the value of each month with the previous one

	var (
		user                = ctx.Locals("user").(*models.User)
		stockEvents         []models.StockEvent
		bondsEvents         []models.BondEvent
		stocksAsInvestments []models.Investment
		response            []models.GenericLabelValue
	)

	database.BankDB.Order("created_at").Find(&stockEvents, "user_id = ?", user.ID)
	database.BankDB.Order("created_at").Find(&bondsEvents, "user_id = ?", user.ID)

	for _, stk := range stockEvents {
		stocksAsInvestments = append(stocksAsInvestments, models.StockToInvestment(stk))
	}
	for _, bds := range bondsEvents {
		stocksAsInvestments = append(stocksAsInvestments, models.BondToInvestment(bds))
	}

	monthlyInvestments := groupMonthly(stocksAsInvestments)
	cumulativeMonthlyInvestments := cumulativeMonthly(monthlyInvestments)
	evolution := fillDateGapsForEvolution(cumulativeMonthlyInvestments)

	for _, v := range evolution {
		response = append(response, models.GenericLabelValue{Label: v.CreatedAt.Format("Jan-06"), Value: uint(v.TotalPrice)})
	}

	return ctx.JSON(response)
}

func groupMonthly(invts []models.Investment) []models.Investment {

	investmentsMap := make(map[time.Time]models.Investment)
	var monthly_contribution []models.Investment

	for _, ivt := range invts {
		year, month, _ := ivt.CreatedAt.UTC().Date()
		finalDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)

		if monthly_ivt, ok := investmentsMap[finalDate]; ok == true {
			monthly_ivt.TotalPrice = monthly_ivt.TotalPrice + ivt.TotalPrice
			monthly_ivt.Quantity = monthly_ivt.Quantity + ivt.Quantity
			monthly_ivt.AveragePrice = monthly_ivt.TotalPrice / float64(monthly_ivt.Quantity)
			investmentsMap[finalDate] = monthly_ivt
		} else {
			ivt.CreatedAt = finalDate
			investmentsMap[finalDate] = ivt
		}
	}

	for _, v := range investmentsMap {
		monthly_contribution = append(monthly_contribution, v)
	}

	sort.Slice(monthly_contribution, func(i, j int) bool {
		return monthly_contribution[i].CreatedAt.Before(monthly_contribution[j].CreatedAt)
	})

	return monthly_contribution
}

func cumulativeMonthly(monthly_contribution []models.Investment) (cumulative_monthly []models.Investment) {

	sort.Slice(monthly_contribution, func(i, j int) bool {
		return monthly_contribution[i].CreatedAt.Before(monthly_contribution[j].CreatedAt)
	})

	cumulative_monthly = make([]models.Investment, len(monthly_contribution))
	cumulative_monthly[0] = monthly_contribution[0]

	for i := 1; i < len(monthly_contribution); i++ {

		qty := monthly_contribution[i].Quantity + cumulative_monthly[i-1].Quantity
		tprice := monthly_contribution[i].TotalPrice + cumulative_monthly[i-1].TotalPrice
		cumulative_monthly[i] = models.Investment{
			UserID:       monthly_contribution[i].UserID,
			StockID:      monthly_contribution[i].StockID,
			CreatedAt:    monthly_contribution[i].CreatedAt,
			Quantity:     qty,
			TotalPrice:   tprice,
			AveragePrice: tprice / float64(qty),
		}
	}
	return
}

func fillDateGapsForEvolution(data []models.Investment) []models.Investment {
	if len(data) <= 1 {
		return data
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].CreatedAt.Before(data[j].CreatedAt)
	})

	y, m, _ := time.Now().Date()
	currentDate := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
	data = append(data, models.Investment{CreatedAt: currentDate})
	var missingData []models.Investment

	for i := 0; i < len(data)-1; i++ {
		y0, m0, _ := data[i].CreatedAt.Date()
		date0 := time.Date(y0, m0, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
		y1, m1, _ := data[i+1].CreatedAt.Date()
		date1 := time.Date(y1, m1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

		for date1.After(date0) {
			y0, m0, _ := date0.Date()
			date0 = time.Date(y0, m0, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 2, -1)
			missingData = append(missingData, models.Investment{
				UserID:       data[i].UserID,
				StockID:      data[i].StockID,
				Quantity:     data[i].Quantity,
				AveragePrice: data[i].AveragePrice,
				TotalPrice:   data[i].TotalPrice,
				CreatedAt:    date0,
			})
		}
	}

	data = data[:len(data)-1]
	data = append(data, missingData...)

	sort.Slice(data, func(i, j int) bool {
		return data[i].CreatedAt.Before(data[j].CreatedAt)
	})

	return data
}

// Simulations
func GetInvestmentsEvolutionSimulation(ctx *fiber.Ctx) error {
	// For the selected index
	// Simulate the cumulative previous month value + index of the current month

	var (
		user                    = ctx.Locals("user").(*models.User)
		statements              []models.Statement
		statementsAsInvestments []models.Investment
	)

	database.BankDB.Order("created_at").Find(&statements, "event = 'TRANSFER'", "user_id = ?", user.ID)
	// TODO make a return when statements is empty

	for _, stmt := range statements {
		statementsAsInvestments = append(statementsAsInvestments, models.StatementToInvestment(stmt))
	}

	monthlyInvestments := groupMonthly(statementsAsInvestments)
	groupedInvestmentsByMonth := cumulativeMonthly(monthlyInvestments)
	cumulativeMonthlyInvestments := fillDateGapsForEvolution(groupedInvestmentsByMonth)
	simulatedMonthlyEv := simulateEvolution("ipca", monthlyInvestments)

	var response models.Simulation
	var values []models.MonthlyValues
	response.Lines = append(response.Lines, "ipca", "deposits")
	for i, v := range simulatedMonthlyEv {
		values = append(values, models.MonthlyValues{Date: v.CreatedAt.Format("Jan-06")})
		values[len(values)-1].Value = append(values[len(values)-1].Value, uint(v.TotalPrice), uint(cumulativeMonthlyInvestments[i].TotalPrice))
	}
	response.Values = values
	return ctx.JSON(response)
}

func simulateEvolution(indexer string, monthly_contribution []models.Investment) (simulatedMonthlyEv []models.Investment) {
	var (
		ipcaHistory  []models.IPCA
		selicHistory []models.SELIC
	)

	sort.Slice(monthly_contribution, func(i, j int) bool {
		return monthly_contribution[i].CreatedAt.Before(monthly_contribution[j].CreatedAt)
	})

	startingDate := monthly_contribution[0].CreatedAt.UTC().AddDate(0, 0, -1)
	database.BankDB.Order("date").Find(&ipcaHistory, "date > ?", startingDate)
	database.BankDB.Order("date").Find(&selicHistory, "date > ?", startingDate)

	montly_contribution_map := make(map[time.Time]models.Investment)
	for _, invt := range monthly_contribution {
		montly_contribution_map[invt.CreatedAt] = invt
	}

	// Get a month of investment multiple by the indexer on the imediatly following month, store the result
	// Get the previous result, add it to the imediatly following month of investiments, apply the result in the previous step
	// EQ => Ms[i] := { Ms[i-1] x (1+IDX[i-1]) } + M[i]

	simulatedInvt := make(map[time.Time]models.Investment)
	simulatedInvt[selicHistory[0].Date.UTC()] = models.Investment{
		TotalPrice: montly_contribution_map[selicHistory[0].Date.UTC()].TotalPrice,
		CreatedAt:  selicHistory[0].Date.UTC(),
	}
	for i := 1; i < len(selicHistory); i++ {
		year, month, _ := selicHistory[i].Date.UTC().Date()
		currentMonth := selicHistory[i].Date.UTC()
		previousMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
		simulatedInvt[currentMonth] = models.Investment{
			TotalPrice: (simulatedInvt[previousMonth].TotalPrice * (1 + selicHistory[i-1].Rate/100)) + montly_contribution_map[currentMonth].TotalPrice,
			CreatedAt:  currentMonth,
		}
	}

	for _, v := range simulatedInvt {
		simulatedMonthlyEv = append(simulatedMonthlyEv, v)
	}

	sort.Slice(simulatedMonthlyEv, func(i, j int) bool {
		return simulatedMonthlyEv[i].CreatedAt.Before(simulatedMonthlyEv[j].CreatedAt)
	})

	return
}
