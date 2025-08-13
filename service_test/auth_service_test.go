package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"testlake/inout/auth"

	"github.com/joho/godotenv"
)

func TestAuthService(t *testing.T) {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// Test signup
	t.Run("SignUp", func(t *testing.T) {
		err, response := SignUpHelper()
		if err != nil {
			t.Fatalf("Failed to sign up: %v", err)
		}
		fmt.Println("SignUp successful:", response.Data.User.Email)
	})

	// Test signin
	t.Run("SignIn", func(t *testing.T) {
		err, response := SignInHelper()
		if err != nil {
			t.Fatalf("Failed to sign in: %v", err)
		}
		fmt.Println("SignIn successful:", response.Data.User.Email)
	})

	// Test refresh token
	t.Run("RefreshToken", func(t *testing.T) {
		// First signin to get a token
		err, signinResponse := SignInHelper()
		if err != nil {
			t.Fatalf("Failed to sign in for refresh token test: %v", err)
		}

		err, response := RefreshTokenHelper(signinResponse.Data.Token)
		if err != nil {
			t.Fatalf("Failed to refresh token: %v", err)
		}
		fmt.Println("RefreshToken successful, new token length:", len(response.Data.Token))
	})

	// Test forgot password
	t.Run("ForgotPassword", func(t *testing.T) {
		err := ForgotPasswordHelper()
		if err != nil {
			t.Fatalf("Failed to request password reset: %v", err)
		}
		fmt.Println("ForgotPassword request successful")
	})
}

func SignUpHelper() (error, auth.SignUpOut) {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/auth/signup"

	// Generate unique email for testing
	userData := auth.SignUpRequest{
		Email:        fmt.Sprintf("testuser%d@example.com", getUniqueID()),
		Username:     fmt.Sprintf("testuser%d", getUniqueID()),
		Password:     "password123",
		FirstName:    stringPtr("Test"),
		LastName:     stringPtr("User"),
		AuthProvider: "email",
	}

	jsonData, _ := json.Marshal(userData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err, auth.SignUpOut{}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, auth.SignUpOut{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, auth.SignUpOut{}
	}

	var result auth.SignUpOut
	if err := json.Unmarshal(body, &result); err != nil {
		return err, auth.SignUpOut{}
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Expected status 201, got %v, body: %s", resp.StatusCode, string(body)), auth.SignUpOut{}
	}

	return nil, result
}

func SignInHelper() (error, auth.SignInOut) {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/auth/signin"

	// First create a user
	signupErr, signupResponse := SignUpHelper()
	if signupErr != nil {
		return signupErr, auth.SignInOut{}
	}

	signinData := auth.SignInRequest{
		Email:    signupResponse.Data.User.Email,
		Password: "password123",
	}

	jsonData, _ := json.Marshal(signinData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err, auth.SignInOut{}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, auth.SignInOut{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, auth.SignInOut{}
	}

	var result auth.SignInOut
	if err := json.Unmarshal(body, &result); err != nil {
		return err, auth.SignInOut{}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200, got %v, body: %s", resp.StatusCode, string(body)), auth.SignInOut{}
	}

	return nil, result
}

func RefreshTokenHelper(token string) (error, auth.RefreshTokenOut) {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/auth/refresh"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err, auth.RefreshTokenOut{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, auth.RefreshTokenOut{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, auth.RefreshTokenOut{}
	}

	var result auth.RefreshTokenOut
	if err := json.Unmarshal(body, &result); err != nil {
		return err, auth.RefreshTokenOut{}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200, got %v, body: %s", resp.StatusCode, string(body)), auth.RefreshTokenOut{}
	}

	return nil, result
}

func ForgotPasswordHelper() error {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/auth/forgot-password"

	forgotPasswordData := auth.ForgotPasswordRequest{
		Email: "test@example.com",
	}

	jsonData, _ := json.Marshal(forgotPasswordData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Expected status 200, got %v, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Helper functions
func getUniqueID() int64 {
	return time.Now().UnixNano()
}

func stringPtr(s string) *string {
	return &s
}
