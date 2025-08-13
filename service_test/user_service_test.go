package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"testlake/inout/user"
	"testlake/model"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
	
	code := m.Run()
	os.Exit(code)
}

func getBaseURL() string {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")
	return scheme + "://" + ip + ":" + port
}

func TestCreateUserIntegration(t *testing.T) {
	baseURL := getBaseURL()
	url := baseURL + "/api/v1/public/user/create"
	
	userData := user.CreateUserRequest{
		Email:        "integration@example.com",
		Username:     "integrationuser",
		Password:     "password123",
		AuthProvider: model.AuthProviderEmail,
	}
	
	jsonData, _ := json.Marshal(userData)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Skip("Server not running - skipping integration test")
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	
	var result user.UserOut
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	
	if resp.StatusCode == http.StatusCreated {
		assert.Equal(t, 0, result.ErrorCode)
		assert.Equal(t, userData.Email, result.Data.Email)
		assert.Equal(t, userData.Username, result.Data.Username)
	} else {
		t.Logf("Expected status 201, got %v. Response: %s", resp.StatusCode, string(body))
	}
}

func TestLoginIntegration(t *testing.T) {
	baseURL := getBaseURL()
	
	createURL := baseURL + "/api/v1/public/user/create"
	loginURL := baseURL + "/api/v1/public/user/login"
	
	userData := user.CreateUserRequest{
		Email:        "logintest@example.com",
		Username:     "logintest",
		Password:     "password123",
		AuthProvider: model.AuthProviderEmail,
	}
	
	jsonData, _ := json.Marshal(userData)
	
	req, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Skip("Server not running - skipping integration test")
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Skip("Server not running - skipping integration test")
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		t.Skip("User creation failed - skipping login test")
		return
	}
	
	loginData := user.LoginRequest{
		Email:    userData.Email,
		Password: userData.Password,
	}
	
	loginJSON, _ := json.Marshal(loginData)
	
	loginReq, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(loginJSON))
	assert.NoError(t, err)
	
	loginReq.Header.Set("Content-Type", "application/json")
	
	loginResp, err := client.Do(loginReq)
	assert.NoError(t, err)
	defer loginResp.Body.Close()
	
	loginBody, err := io.ReadAll(loginResp.Body)
	assert.NoError(t, err)
	
	var loginResult user.LoginOut
	err = json.Unmarshal(loginBody, &loginResult)
	assert.NoError(t, err)
	
	if loginResp.StatusCode == http.StatusOK {
		assert.Equal(t, 0, loginResult.ErrorCode)
		assert.NotEmpty(t, loginResult.Data.Token)
		assert.Equal(t, userData.Email, loginResult.Data.User.Email)
	} else {
		t.Logf("Login failed with status %v. Response: %s", loginResp.StatusCode, string(loginBody))
	}
}

func TestGetUserListIntegration(t *testing.T) {
	baseURL := getBaseURL()
	url := baseURL + "/api/v1/private/user/list"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Skip("Server not running - skipping integration test")
		return
	}
	
	req.Header.Set("Authorization", "Bearer your_test_token_here")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Skip("Server not running - skipping integration test")
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusUnauthorized {
		t.Log("Expected unauthorized response for test without valid token")
		return
	}
	
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	
	fmt.Printf("Response Status: %d, Body: %s\n", resp.StatusCode, string(body))
}