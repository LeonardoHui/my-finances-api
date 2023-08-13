package utils

import (
	"bufio"
	"log"
	"my-finances-api/src/database"
	"my-finances-api/src/models"
	"os"
	"time"
)

// For test only
func InternalCreateNewUser() {
	user := models.User{
		Name:     "test",
		Username: "test",
		Password: "test",
		Email:    "test@email.com",
	}
	if err := user.HashPassword(user.Password); err != nil {
		log.Println("Error hashing password", err)
	}
	if err := database.BankDB.Create(&user).Error; err != nil {
		log.Println("Error saving DB", err)
	}
}

func InternalLoadTables(dirPath string) {

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Println("Error reading internal dir. ", err)
		return
	}
	for _, file := range files {
		readFile, err := os.Open(dirPath + "/" + file.Name())
		if err != nil {
			log.Println("Error opening file. ", err)
		}
		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			database.BankDB.Exec(fileScanner.Text())
		}
		readFile.Close()
	}
}

type Timeline struct {
	Date  time.Time
	Value uint
}

func MonthsArray(startDate time.Time) (timeline []Timeline) {
	y, m, _ := time.Now().Date()
	currentDate := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	for currentDate.After(startDate) {
		year, month, _ := startDate.Date()
		lastDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
		timeline = append(timeline, Timeline{Date: lastDayOfMonth})

		startDate = lastDayOfMonth.AddDate(0, 0, 1)
	}
	return
}
