package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yunusemreayhan/goAuthMicroService/models"
	"github.com/yunusemreayhan/goAuthMicroService/restapi"
)

const (
	testUsername  = "test3"
	testPassword  = "test3"
	testPassword2 = "test3_2"
	testEmail     = "test3@gmail.com"
)

var client = &http.Client{
	Transport: &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	},
	Timeout: 30 * time.Second,
}

var server *httptest.Server

var httpHandler *http.Handler

func bodyCloser(body io.ReadCloser) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Default().Printf("Body.Close error : [%v]\n", err)
		}
	}(body)
}
func TestAuthApi(t *testing.T) {
	if httpHandler == nil {
		handler, err := restapi.GetAPIHandler()
		if err != nil {
			t.Fatal("get api handler", err)
		}
		httpHandler = &handler
	}
	if server == nil {
		server = httptest.NewServer(*httpHandler)
	}

	registrationResponse := registerTestPerson(t)

	var err error
	_, err = login(t, err, testUsername, testPassword+"wrong")
	require.Error(t, err)

	token, err := login(t, err, testUsername, testPassword)
	require.NoError(t, err)

	err = verifyToken(t, err, token)
	require.NoError(t, err)

	err = verifyToken(t, err, "wrong")
	require.Error(t, err)

	err = changePassword(t, err, registrationResponse.ID, token)

	_, err = login(t, err, testUsername, testPassword)
	require.Error(t, err)
	token, err = login(t, err, testUsername, testPassword2)
	require.NoError(t, err)

	deletePerson(t, err, registrationResponse.ID, token)

	registrationResponse = registerTestPerson(t)
}

func verifyToken(t *testing.T, err error, token string) error {
	verifyReq, err := http.NewRequest("GET", server.URL+"/api/verify", strings.NewReader(""))
	require.NoError(t, err)
	verifyReq.Header.Set("Content-Type", "application/json")
	verifyReq.Header.Add("oauth_token", token)
	verifyResp, err := client.Do(verifyReq)
	require.NoError(t, err)
	if verifyResp.StatusCode != http.StatusOK {
		var mError models.Error
		err = json.NewDecoder(verifyResp.Body).Decode(&mError)
		defer bodyCloser(verifyResp.Body)
		require.NoError(t, err)
		log.Default().Printf("verify token resp : [%#v]\n", mError)
		return fmt.Errorf("verify token resp : [%#v]", mError)
	}
	return nil
}

func changePassword(t *testing.T, err error, id int64, token string) error {
	changePassReq, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/person/%d", server.URL, id), strings.NewReader(fmt.Sprintf("{\"person_name\":\"%s\",\"email\":\"%s\",\"password\":\"%s\"}", testUsername, testEmail, testPassword2)))
	require.NoError(t, err)
	changePassReq.Header.Set("Content-Type", "application/json")
	changePassReq.Header.Add("oauth_token", token)
	changePassResp, err := client.Do(changePassReq)
	require.NoError(t, err)
	if changePassResp.StatusCode != http.StatusOK {
		var mError models.Error
		err = json.NewDecoder(changePassResp.Body).Decode(&mError)
		defer bodyCloser(changePassResp.Body)
		require.NoError(t, err)
		log.Default().Printf("change password resp : [%#v]\n", mError)
	}
	require.Equal(t, http.StatusOK, changePassResp.StatusCode)
	return err
}

func deletePerson(t *testing.T, err error, id int64, token string) {
	deleteReq, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/person/%d", server.URL, id), strings.NewReader(""))
	require.NoError(t, err)
	deleteReq.Header.Set("Content-Type", "application/json")
	deleteReq.Header.Add("oauth_token", token)
	deleteResp, err := client.Do(deleteReq)
	require.NoError(t, err)
	if deleteResp.StatusCode != http.StatusNoContent {
		var mError models.Error
		err = json.NewDecoder(deleteResp.Body).Decode(&mError)
		defer bodyCloser(deleteResp.Body)
		require.NoError(t, err)
		log.Default().Printf("delete person resp : [%#v]\n", mError)
	}
	require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)
}

func registerTestPerson(t *testing.T) models.RegistrationResponse {
	// register
	registerReq, err := http.NewRequest("POST", server.URL+"/api/register", strings.NewReader(fmt.Sprintf(`{"person_name":"%s","email":"%s","password":"%s"}`, testUsername, testEmail, testPassword)))
	require.NoError(t, err)
	registerReq.Header.Set("Content-Type", "application/json")
	respFromRegister, err := client.Do(registerReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, respFromRegister.StatusCode)
	var registrationResponse models.RegistrationResponse
	err = json.NewDecoder(respFromRegister.Body).Decode(&registrationResponse)
	defer bodyCloser(respFromRegister.Body)
	require.NoError(t, err)
	log.Default().Printf("resp : [%#v]\n", registrationResponse)
	return registrationResponse
}

func login(t *testing.T, err error, username string, password string) (string, error) {
	loginReq, err := http.NewRequest("POST", server.URL+"/api/login", strings.NewReader(""))
	require.NoError(t, err)
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.SetBasicAuth(username, password)
	require.NoError(t, err)
	respFromLogin, err := client.Do(loginReq)
	if err != nil || respFromLogin.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login error : [%v] [%v]", err, respFromLogin.StatusCode)
	}
	require.NoError(t, err)
	var lResp models.LoginResponse
	err = json.NewDecoder(respFromLogin.Body).Decode(&lResp)
	defer bodyCloser(respFromLogin.Body)
	log.Default().Printf("resp : [%#v]\n", lResp)
	return *lResp.Token, nil
}
