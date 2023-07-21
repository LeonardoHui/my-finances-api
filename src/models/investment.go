package models

import (
	"time"

	"gorm.io/gorm"
)

type Investment struct {
	UserID       uint
	StockID      uint
	Quantity     uint
	AveragePrice float64
	TotalPrice   float64
	CreatedAt    time.Time
}

type InvestmentEvent struct {
	gorm.Model
	UserID      uint
	Amount      int64
	Description string
	StockID     uint
}

func BondToInvestment(be BondEvent) Investment {
	return Investment{
		UserID:     be.UserID,
		StockID:    be.BondID,
		Quantity:   be.Quantity,
		TotalPrice: float64(be.TotalPrice),
		CreatedAt:  be.CreatedAt,
	}
}

func StockToInvestment(se StockEvent) Investment {
	return Investment{
		UserID:       se.UserID,
		StockID:      se.StockID,
		Quantity:     se.Quantity,
		AveragePrice: float64(se.Price),
		TotalPrice:   float64(se.TotalPrice),
		CreatedAt:    se.CreatedAt,
	}
}

func StatementToInvestment(stmt Statement) Investment {
	return Investment{
		UserID:       stmt.UserID,
		StockID:      0,
		Quantity:     1,
		AveragePrice: 0,
		TotalPrice:   stmt.Amount,
		CreatedAt:    stmt.CreatedAt,
	}
}
