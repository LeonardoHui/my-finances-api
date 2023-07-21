package models

import "gorm.io/gorm"

type Broker struct {
	gorm.Model
	Name string
}
