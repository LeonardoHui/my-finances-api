package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var BankDB *gorm.DB

type DbConfigs struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SslMode  string
	TimeZone string
}

// Example: dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
func (db DbConfigs) dns() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		db.Host, db.User, db.Password, db.Name, db.Port,
	)
}

func (db DbConfigs) Open() *gorm.DB {
	dns := db.dns()
	session, err := gorm.Open(postgres.Open(dns))
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return session
}
