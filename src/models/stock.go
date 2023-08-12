package models

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	Name string
}

type StockEvent struct {
	gorm.Model
	UserID          uint
	BrokerID        uint
	StockID         uint
	Event           string // BUY | SELL | DIVIDEND
	Quantity        uint
	Price           uint //In Cents
	LiquidationFee  uint //In Cents
	RegistrationFee uint //In Cents
	NegociationFee  uint //In Cents
	TotalPrice      uint //In Cents
}

type StockHistory struct {
	gorm.Model
	StockID uint
	Price   uint //In Cents
}

func (StockHistory) TableName() string {
	return "stock_history"
}
