package models

import "gorm.io/gorm"

type Bank struct {
	gorm.Model
	Name string
}
