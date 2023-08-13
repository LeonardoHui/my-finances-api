package models

import "time"

type GenericMonetaryItem struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Amount      uint      `json:"amount"`
	Date        time.Time `json:"date"`
}

type UserStatements struct {
	Statements []GenericMonetaryItem `json:"statements"`
	Balance    []GenericMonetaryItem `json:"balance"`
}

type GenericLabelValue struct {
	Label string `json:"label"`
	Value uint   `json:"value"`
}

type UserInvestments struct {
	Investments   []GenericMonetaryItem `json:"investments"`
	DividendYield []GenericLabelValue   `json:"dividend_yied"`
	DividendPaid  []GenericLabelValue   `json:"dividend_paid"`
}

type MonthlyValues struct {
	Date  string `json:"date"`
	Value []uint `json:"value"`
}

type Simulation struct {
	Lines  []string        `json:"lines"`
	Values []MonthlyValues `json:"values"`
}
