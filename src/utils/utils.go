package utils

import (
	"bufio"
	"log"
	"my-finances-api/src/database"
	"my-finances-api/src/models"
	"os"
)

// For teste only
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
		readFile, err := os.Open(dirPath + file.Name())
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
