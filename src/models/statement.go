package models

import "gorm.io/gorm"

type Statement struct {
	gorm.Model
	UserID      uint
	Event       string  // TRANSFER | BUY | SELL | DIVIDEND | PAYMENT | LIQUIDATION
	Amount      float64 // Positve is IN | Negative is OUT
	Description string
	BankID      uint
}
