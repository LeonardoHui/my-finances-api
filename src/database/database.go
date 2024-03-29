package database

import (
	"fmt"
	"log"
	"my-finances-api/src/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

func OpenLite(fileName string) *gorm.DB {
	session, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return session
}

func Migrate() {
	BankDB.AutoMigrate(models.User{})
	BankDB.AutoMigrate(models.Bank{})
	BankDB.AutoMigrate(models.BankAccount{})
	BankDB.AutoMigrate(models.Statement{})
	BankDB.AutoMigrate(models.Stock{})
	BankDB.AutoMigrate(models.StockEvent{})
	BankDB.AutoMigrate(models.StockHistory{})
	BankDB.AutoMigrate(models.Bond{})
	BankDB.AutoMigrate(models.BondEvent{})
	BankDB.AutoMigrate(models.BondHistory{})
	BankDB.AutoMigrate(models.IPCA{})
	BankDB.AutoMigrate(models.SELIC{})
}
