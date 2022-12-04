package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_integration(t *testing.T) {

	url := "http://127.0.0.1:3000/stock/MXRF11/price"
	req, _ := http.NewRequest("GET", url, nil)

	resp, _ := http.DefaultClient.Do(req)
	resBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resBody))
}
