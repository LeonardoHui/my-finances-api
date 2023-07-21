package models

import "time"

// Specific configuration for table Ipca
type IPCA struct {
	Date time.Time `gorm:"primarykey"`
	Rate float64
}

func (IPCA) TableName() string {
	return "ipca"
}

// Specific configuration for table Selic
type SELIC struct {
	Date time.Time `gorm:"primarykey"`
	Rate float64
}

func (SELIC) TableName() string {
	return "selic"
}
