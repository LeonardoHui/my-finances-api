package models

import "gorm.io/gorm"

type Statement struct {
	gorm.Model
	UserID      uint
	Amount      float64
	Description string
	BankID      uint
}
