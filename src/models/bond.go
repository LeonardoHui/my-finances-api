package models

import (
	"time"

	"gorm.io/gorm"
)

type Bond struct {
	gorm.Model
	Name string
}

type BondEvent struct {
	gorm.Model
	UserID     uint
	BrokerID   uint
	BondID     uint
	Event      string // BUY | SELL | FEE | PAYOUT
	Quantity   uint
	Price      uint //In Cents
	TotalPrice uint //In Cents
	Index      string
	Rate       uint
	Maturity   time.Time
}

type BondHistory struct {
	gorm.Model
	StockID uint
	Price   uint //In Cents
}

func (BondHistory) TableName() string {
	return "bond_history"
}
