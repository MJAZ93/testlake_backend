package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"testlake/app"
	"testlake/dao"
	"testlake/inout/auth"
	_ "testlake/inout/user"
	"testlake/model"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}

	dao.Connect()
	gin.SetMode(gin.TestMode)

	code := m.Run()
	os.Exit(code)
}

func setupRouter() *gin.Engine {
	router := gin.New()

	baseRoute := router.Group("/api/v1")

	publicRoutes := baseRoute.Group("/public")
	app.PublicRoutes(publicRoutes)

	return router
}

func TestCreateUser(t *testing.T) {
	router := setupRouter()

	uniqueID := time.Now().UnixNano()
	requestData := auth.SignUpRequest{
		Email:        fmt.Sprintf("testcontroller%d@example.com", uniqueID),
		Username:     fmt.Sprintf("testcontroller%d", uniqueID),
		Password:     "password123",
		AuthProvider: model.AuthProviderEmail,
	}

	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/api/v1/public/auth/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response auth.SignUpOut
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.ErrorCode)
	assert.Equal(t, requestData.Email, response.Data.User.Email)
	assert.Equal(t, requestData.Username, response.Data.User.Username)

	userDao := dao.NewUserDao()
	defer userDao.Delete(response.Data.User.ID)
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	router := setupRouter()
	userDao := dao.NewUserDao()

	uniqueID := time.Now().UnixNano()
	duplicateEmail := fmt.Sprintf("duplicate%d@example.com", uniqueID)

	existingUser := &model.User{
		Email:        duplicateEmail,
		Username:     fmt.Sprintf("existing%d", uniqueID),
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}

	err := userDao.Create(existingUser)
	assert.NoError(t, err)
	defer userDao.Delete(existingUser.ID)

	requestData := auth.SignUpRequest{
		Email:        duplicateEmail,
		Username:     fmt.Sprintf("newuser%d", uniqueID),
		Password:     "password123",
		AuthProvider: model.AuthProviderEmail,
	}

	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/api/v1/public/auth/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin(t *testing.T) {
	router := setupRouter()
	userDao := dao.NewUserDao()

	uniqueID := time.Now().UnixNano()
	password := "password123"
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	hashedPassword := string(hashedPasswordBytes)

	email := fmt.Sprintf("login%d@example.com", uniqueID)
	testUser := &model.User{
		Email:        email,
		Username:     fmt.Sprintf("loginuser%d", uniqueID),
		AuthProvider: model.AuthProviderEmail,
		PasswordHash: &hashedPassword,
		Status:       model.UserStatusActive,
	}

	err = userDao.Create(testUser)
	assert.NoError(t, err)
	defer userDao.Delete(testUser.ID)

	loginRequest := auth.SignInRequest{
		Email:    email,
		Password: password,
	}

	jsonData, _ := json.Marshal(loginRequest)

	req, _ := http.NewRequest("POST", "/api/v1/public/auth/signin", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Test should expect successful login
	assert.Equal(t, http.StatusOK, w.Code)

	var response auth.SignInOut
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.ErrorCode)
	assert.NotEmpty(t, response.Data.Token)
	assert.Equal(t, email, response.Data.User.Email)
}
