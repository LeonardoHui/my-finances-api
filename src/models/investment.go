package models

import "gorm.io/gorm"

type Investiment struct {
	gorm.Model
	UserID       uint
	StockID      uint
	Quantity     uint
	AveragePrice float64
	Total        float64
}

type InvestimentEvent struct {
	gorm.Model
	UserID      uint
	Amount      int64
	Description string
	StockID     uint
}
