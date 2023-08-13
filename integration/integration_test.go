package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"my-finances-api/src/configs"
	"my-finances-api/src/database"
	"my-finances-api/src/models"
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

	configs.Envs["ENV"] = "TEST"
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

func Test_Statements_Endpoint_Response_Shall_Not_Be_Empty_With_StatusOK(t *testing.T) {

	url := host + "/statements"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	responseFormat := models.UserStatements{}

	if err := json.Unmarshal(respBody, &responseFormat); err != nil {
		t.Log("Unable to unmarshal response body")
		t.Fail()
	}
	assert.NotEmpty(t, responseFormat, "Response body was empty")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Investments_Endpoint_Response_Shall_Not_Be_Empty_With_StatusOK(t *testing.T) {

	url := host + "/investments"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	responseFormat := models.UserInvestments{}

	if err := json.Unmarshal(respBody, &responseFormat); err != nil {
		t.Log("Unable to unmarshal response body")
		t.Fail()
	}
	assert.NotEmpty(t, responseFormat, "Response body was empty")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Evolution_Endpoint_Response_Shall_Not_Be_Empty_With_StatusOK(t *testing.T) {

	url := host + "/evolution"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	responseFormat := []models.GenericLabelValue{}

	if err := json.Unmarshal(respBody, &responseFormat); err != nil {
		t.Log("Unable to unmarshal response body")
		t.Fail()
	}
	assert.NotEmpty(t, responseFormat, "Response body was empty")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Simulation_Endpoint_Response_Shall_Not_Be_Empty_With_StatusOK(t *testing.T) {

	url := host + "/simulation"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	responseFormat := models.Simulation{}

	if err := json.Unmarshal(respBody, &responseFormat); err != nil {
		t.Log("Unable to unmarshal response body")
		t.Fail()
	}
	assert.NotEmpty(t, responseFormat, "Response body was empty")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
