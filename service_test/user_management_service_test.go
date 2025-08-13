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

	"github.com/joho/godotenv"
)

func TestUserManagementService(t *testing.T) {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// First create a user and get auth token
	err, authResponse := SignUpHelper()
	if err != nil {
		t.Fatalf("Failed to sign up for user management tests: %v", err)
	}

	token := authResponse.Data.Token

	// Test get profile
	t.Run("GetProfile", func(t *testing.T) {
		err, response := GetProfileHelper(token)
		if err != nil {
			t.Fatalf("Failed to get profile: %v", err)
		}
		fmt.Println("GetProfile successful:", response.Data.Email)
	})

	// Test update profile
	t.Run("UpdateProfile", func(t *testing.T) {
		err, response := UpdateProfileHelper(token)
		if err != nil {
			t.Fatalf("Failed to update profile: %v", err)
		}
		fmt.Println("UpdateProfile successful:", *response.Data.FirstName)
	})

	// Test get dashboard
	t.Run("GetDashboard", func(t *testing.T) {
		err, response := GetDashboardHelper(token)
		if err != nil {
			t.Fatalf("Failed to get dashboard: %v", err)
		}
		fmt.Println("GetDashboard successful, projects:", response.Data.PersonalProjects)
	})

	// Test get notifications
	t.Run("GetNotifications", func(t *testing.T) {
		err, response := GetNotificationsHelper(token)
		if err != nil {
			t.Fatalf("Failed to get notifications: %v", err)
		}
		fmt.Println("GetNotifications successful, count:", len(response.Data))
	})

	// Test delete account (should be last)
	t.Run("DeleteAccount", func(t *testing.T) {
		err := DeleteAccountHelper(token)
		if err != nil {
			t.Fatalf("Failed to delete account: %v", err)
		}
		fmt.Println("DeleteAccount successful")
	})
}

func GetProfileHelper(token string) (error, user.UserOut) {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/users/profile"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, user.UserOut{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, user.UserOut{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, user.UserOut{}
	}

	var result user.UserOut
	if err := json.Unmarshal(body, &result); err != nil {
		return err, user.UserOut{}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200, got %v, body: %s", resp.StatusCode, string(body)), user.UserOut{}
	}

	return nil, result
}

func UpdateProfileHelper(token string) (error, user.UserOut) {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/users/profile"

	updateData := user.UpdateUserRequest{
		FirstName: stringPtr("UpdatedFirst"),
		LastName:  stringPtr("UpdatedLast"),
		AvatarURL: stringPtr("https://example.com/avatar.jpg"),
	}

	jsonData, _ := json.Marshal(updateData)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err, user.UserOut{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, user.UserOut{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, user.UserOut{}
	}

	var result user.UserOut
	if err := json.Unmarshal(body, &result); err != nil {
		return err, user.UserOut{}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200, got %v, body: %s", resp.StatusCode, string(body)), user.UserOut{}
	}

	return nil, result
}

func GetDashboardHelper(token string) (error, user.DashboardOut) {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/users/dashboard"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, user.DashboardOut{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, user.DashboardOut{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, user.DashboardOut{}
	}

	var result user.DashboardOut
	if err := json.Unmarshal(body, &result); err != nil {
		return err, user.DashboardOut{}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200, got %v, body: %s", resp.StatusCode, string(body)), user.DashboardOut{}
	}

	return nil, result
}

func GetNotificationsHelper(token string) (error, user.NotificationsOut) {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/users/notifications"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, user.NotificationsOut{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, user.NotificationsOut{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, user.NotificationsOut{}
	}

	var result user.NotificationsOut
	if err := json.Unmarshal(body, &result); err != nil {
		return err, user.NotificationsOut{}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200, got %v, body: %s", resp.StatusCode, string(body)), user.NotificationsOut{}
	}

	return nil, result
}

func DeleteAccountHelper(token string) error {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")

	url := scheme + "://" + ip + ":" + port + "/api/v1/users/account"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

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
