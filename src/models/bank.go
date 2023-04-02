package models

import "gorm.io/gorm"

type Bank struct {
	gorm.Model
	Name string
}

type BankAccount struct {
	gorm.Model
	UserID uint
	Amount int64
	BankID uint
}
