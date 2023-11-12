package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yunusemreayhan/goAuthMicroService/internal/model"
)

var (
	testusername = "test"
	testpassword = "test"
	testemail    = "test@gmail.com"
	voucher      = ""
)

type Data struct {
	url                  string
	json_request         string
	expected_status_code int32
}

func handlerForRegister(datai interface{}) bool {
	data := datai.(Data)
	url := data.url
	jsonData := data.json_request

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonData))
	if err != nil {
		log.Default().Println("Error creating request:", err)
		return false
	}

	// Set the request content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Default().Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == int(data.expected_status_code) {
		return true
	} else {
		log.Default().Println("Expected status code:", data.expected_status_code, "but got:", resp)
		log.Default().Println("Data : ", data)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Default().Println("Error reading response body:", err)
			return false
		}
		log.Default().Println("Response body:", string(bodyBytes))
		return false
	}
}

func handlerForLogin(datai interface{}) bool {
	data := datai.(Data)
	url := data.url
	jsonData := data.json_request

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonData))
	if err != nil {
		log.Default().Println("Error creating request:", err)
		return false
	}

	// Set the request content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Default().Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == int(data.expected_status_code) {
		var loginResp model.LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		if err != nil {
			log.Default().Println("Error decoding response:", err)
			return false
		}
		if data.expected_status_code == 200 {
			voucher = loginResp.Token
			return true
		} else {
			log.Default().Println("Error in reponse :", loginResp.Error)
			return true
		}
	} else {
		log.Default().Println("Expected status code:", data.expected_status_code, "but got:", resp)
		log.Default().Println("Data : ", data)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Default().Println("Error reading response body:", err)
			return false
		}
		log.Default().Println("Response body:", string(bodyBytes))
		return false
	}
}

func handlerForVerify(datai interface{}) bool {
	data := datai.(Data)
	url := data.url
	jsonData := data.json_request

	var verifyRequest model.VerifyVoucherRequest

	err := json.Unmarshal([]byte(jsonData), &verifyRequest)

	if err != nil {
		log.Default().Println("Error decoding request for test:", err)
		return false
	}

	if verifyRequest.Token == "proper" {
		verifyRequest.Token = voucher
	}
	verifyRequestBytes, err := json.Marshal(verifyRequest)
	if err != nil {
		log.Default().Println("Error marshaling verifyRequest:", err)
		return false
	}
	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(verifyRequestBytes)))
	if err != nil {
		log.Default().Println("Error creating request:", err)
		return false
	}

	// Set the request content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Default().Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == int(data.expected_status_code) {
		return true
	} else {
		log.Default().Println("Expected status code:", data.expected_status_code, "but got:", resp)
		log.Default().Println("Data : ", data)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Default().Println("Error reading response body:", err)
			return false
		}
		log.Default().Println("Response body:", string(bodyBytes))
		return false
	}
}

func TestTableForRest(t *testing.T) {
	log.Default().Println("testing table for rest")
	table := []struct {
		data     Data
		testbody func(interface{}) bool
	}{
		{
			Data{
				url:                  "http://localhost/api/register",
				json_request:         "{\"username\":\"test\", \"email\":\"test@test.com\", \"password\":\"12345\"}",
				expected_status_code: 200,
			},
			handlerForRegister,
		},
		{
			Data{
				url:                  "http://localhost/api/login",
				json_request:         "{\"username\":\"test\", \"password\":\"123456\"}",
				expected_status_code: 401,
			},
			handlerForLogin,
		},
		{
			Data{
				url:                  "http://localhost/api/login",
				json_request:         "{\"username\":\"test\", \"password\":\"12345\"}",
				expected_status_code: 200,
			},
			handlerForLogin,
		},
		{
			Data{
				url:                  "http://localhost/api/verify",
				json_request:         "{\"voucher\":\"nonproper\"}",
				expected_status_code: 401,
			},
			handlerForVerify,
		},
		{
			Data{
				url:                  "http://localhost/api/verify",
				json_request:         "{\"voucher\":\"proper\"}",
				expected_status_code: 200,
			},
			handlerForVerify,
		},
	}

	for _, titem := range table {
		log.Default().Println(titem.data.url)
		require.True(t, titem.testbody(titem.data))
	}
}
