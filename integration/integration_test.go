package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"my-finances-api/src/database"
	"my-finances-api/src/server"
	"my-finances-api/src/utils"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type JWT struct {
	Token string `json:"token"`
}

var (
	port = ":8000"
	host = "http://127.0.0.1" + port
	jwt  JWT
)

func TestMain(m *testing.M) {

	log.Printf("Starting tests execution")
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Printf("Error to start tests %v\n", err)
	}

	database.BankDB = db
	database.Migrate()

	utils.InternalCreateNewUser()
	utils.InternalLoadTables("./fixtures")

	go server.Run(port)
	time.Sleep(5 * time.Second)
	exitVal := m.Run()
	log.Printf("Ending tests execution")
	os.Exit(exitVal)
}

func Test_Login(t *testing.T) {

	url := host + "/login"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(`{"email": "test@email.com", "password":"test"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBody, &jwt)
	assert.NotEmpty(t, jwt.Token)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Statements(t *testing.T) {

	url := host + "/statements"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwt.Token)
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.NotEmpty(t, respBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
