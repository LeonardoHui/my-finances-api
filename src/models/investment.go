package models

import "gorm.io/gorm"

type Investment struct {
	gorm.Model
	UserID       uint
	StockID      uint
	Quantity     uint
	AveragePrice float64
	Total        float64
}

type InvestmentEvent struct {
	gorm.Model
	UserID      uint
	Amount      int64
	Description string
	StockID     uint
}
