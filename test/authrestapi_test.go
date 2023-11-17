package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yunusemreayhan/goAuthMicroService/internal/handler"
	"github.com/yunusemreayhan/goAuthMicroService/internal/key"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/yunusemreayhan/goAuthMicroService/internal/model"
)

const (
	testingPort  = 8001
	testUsername = "test3"
	testPassword = "test3"
	testEmail    = "test3@gmail.com"
)

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
)

type Data struct {
	url                string
	jsonRequest        string
	expectedStatusCode int32
}

func registerTestUser() *http.Response {
	regreq := model.RegistrationRequest{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	}
	url := fmt.Sprintf("http://localhost:%d/api/register", testingPort)
	regreqbytes, _ := json.Marshal(regreq)
	log.Default().Println("registering with this request", string(regreqbytes))
	req, err := http.NewRequest("POST", url, strings.NewReader(string(regreqbytes)))
	if err != nil {
		log.Default().Println("Error while sending request", err)
		return nil
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Default().Println("Error sending request:", err)
		return nil
	}
	return resp
}

func getProperVoucher() (*string, error) {
	loginreq := model.LoginRequest{
		Username: testUsername,
		Password: testPassword,
	}
	url := fmt.Sprintf("http://localhost:%d/api/login", testingPort)
	loginreqbytes, _ := json.Marshal(loginreq)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(loginreqbytes)))
	if err != nil {
		log.Default().Println("Error while sending request", err)
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Default().Println("Error sending request:", err)
		return nil, err
	}
	var modelResponse model.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&modelResponse)
	if err != nil {
		log.Default().Println("Error while decoding response", err)
		return nil, err
	}
	if modelResponse.Error != "" {
		log.Default().Println("Error while sending request", modelResponse.Error)
		return nil, fmt.Errorf("error while logging in : [%s]", modelResponse.Error)
	}
	if modelResponse.Token == "" {
		log.Default().Println("Both error and Token is empty!")
		return nil, fmt.Errorf("both error and Token is empty")
	}
	return &modelResponse.Token, nil
}

func handlerForRegister(datai interface{}, client *http.Client) bool {
	data := datai.(Data)
	url := data.url
	jsonData := data.jsonRequest
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonData))

	var restApp handler.RestAPPForAuth
	restApp.PrepareFiberApp()
	restApp.SetPrivateKey(key.GeneratePrivateKey())
	resp, err := restApp.GetFiberApp().Test(req, 30)
	if err != nil {
		log.Default().Println("Error sending request:", err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Default().Println("Error closing response body:", err)
		}
	}(resp.Body)
	if resp.StatusCode == int(data.expectedStatusCode) {
		return true
	} else {
		log.Default().Println("Expected status code:", data.expectedStatusCode, "but got:", resp)
		log.Default().Println("Data : ", data)
		var modelRegisterResponse model.RegistrationResponse
		err = json.NewDecoder(resp.Body).Decode(&modelRegisterResponse)
		if err != nil {
			log.Default().Println("Error reading response body:", err)
			return false
		}
		log.Default().Printf("Response body: [%v]", modelRegisterResponse)
		return false
	}
}

func handlerForLogin(datai interface{}, client *http.Client) bool {
	data := datai.(Data)
	url := data.url
	jsonData := data.jsonRequest

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonData))
	if err != nil {
		log.Default().Println("Error while sending request", err)
		return false
	}
	// Set the request content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Default().Println("Error sending request:", err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Default().Println("Error closing response body:", err)
		}
	}(resp.Body)
	if resp.StatusCode == int(data.expectedStatusCode) {
		var loginResp model.LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		if err != nil {
			log.Default().Println("Error decoding response:", err)
			return false
		}
		if data.expectedStatusCode == 200 {

			return true
		} else {
			log.Default().Println("Error in reponse :", loginResp.Error)
			return true
		}
	} else {
		log.Default().Println("Expected status code:", data.expectedStatusCode, "but got:", resp)
		log.Default().Println("Data : ", data)
		var modelLoginResponse model.LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&modelLoginResponse)
		if err != nil {
			log.Default().Println("Error reading response body:", err)
			return false
		}
		log.Default().Printf("Response body: [%v]", modelLoginResponse)
		return false
	}
}

func handlerForVerify(datai interface{}, client *http.Client) bool {
	data := datai.(Data)
	url := data.url
	jsonData := data.jsonRequest

	var verifyRequest model.VerifyVoucherRequest

	err := json.Unmarshal([]byte(jsonData), &verifyRequest)

	if err != nil {
		log.Default().Println("Error decoding request for test:", err)
		return false
	}

	if verifyRequest.Token == "proper" {
		resp := registerTestUser()
		if resp == nil || resp.StatusCode != 200 {

			var modelRegisterResponse model.RegistrationResponse
			err = json.NewDecoder(resp.Body).Decode(&modelRegisterResponse)
			if err != nil {
				log.Default().Println("Error reading response body:", err)
				return false
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Default().Println("Error closing response body:", err)
				}
			}(resp.Body)
			log.Default().Printf("Error registering test user for voucher [%v] %d", modelRegisterResponse, resp.StatusCode)
			if !strings.Contains(fmt.Sprintf("%v", modelRegisterResponse), "duplicate key value") {
				return false
			}
		}
		token, errGetToken := getProperVoucher()
		if errGetToken != nil {
			log.Default().Println("Error getting proper voucher:", errGetToken)
			return false
		}
		verifyRequest.Token = *token
	}
	verifyRequestBytes, err := json.Marshal(verifyRequest)
	if err != nil {
		log.Default().Println("Error marshaling verifyRequest:", err)
		return false
	}
	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(verifyRequestBytes)))
	if err != nil {
		log.Default().Println("Error while sending request", err)
		return false
	}
	// Set the request content type to JSON
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Default().Println("Error sending request:", err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Default().Println("Error closing response body:", err)
		}
	}(resp.Body)
	if resp.StatusCode == int(data.expectedStatusCode) {
		return true
	} else {
		log.Default().Println("Expected status code:", data.expectedStatusCode, "but got:", resp)
		log.Default().Println("Data : ", data)
		var modelVerifyResponse model.VerifyVoucherResponse
		err = json.NewDecoder(resp.Body).Decode(&modelVerifyResponse)
		if err != nil {
			log.Default().Println("Error reading response body:", err)
			return false
		}
		log.Default().Printf("Response body: [%v]", modelVerifyResponse)
		return false
	}
}

func TestAuthApiRegister(t *testing.T) {
	assert.Equal(t, handlerForRegister(Data{
		url:                fmt.Sprintf("http://localhost:%d/api/register", testingPort),
		jsonRequest:        `{"username":"test", "email":"test@test.com", "password":"12345"}`,
		expectedStatusCode: 200,
	}, httpClient), true)
}

func TestAuthApiLoginProperCredentials(t *testing.T) {
	assert.Equal(t, handlerForLogin(Data{
		url:                fmt.Sprintf("http://localhost:%d/api/login", testingPort),
		jsonRequest:        `{"username":"test", "email":"test@test.com", "password":"12345"}`,
		expectedStatusCode: 200,
	}, httpClient), true)
}

func TestAuthApiLoginWrongCredentials(t *testing.T) {
	assert.Equal(t, handlerForLogin(Data{
		url:                fmt.Sprintf("http://localhost:%d/api/login", testingPort),
		jsonRequest:        `{"username":"test", "password":"123456"}`,
		expectedStatusCode: 401,
	}, httpClient), true)
}

func TestAuthApiVerifyProperToken(t *testing.T) {
	assert.Equal(t, handlerForVerify(Data{
		url:                fmt.Sprintf("http://localhost:%d/api/verify", testingPort),
		jsonRequest:        `{"voucher":"proper"}`,
		expectedStatusCode: 200,
	}, httpClient), true)
}

func TestAuthApiVerifyWrongToken(t *testing.T) {
	assert.Equal(t, handlerForVerify(Data{
		url:                fmt.Sprintf("http://localhost:%d/api/verify", testingPort),
		jsonRequest:        `{"voucher":"nonproper"}`,
		expectedStatusCode: 401,
	}, httpClient), true)
}
