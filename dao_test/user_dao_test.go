package dao_test

import (
	"os"
	"testing"

	"testlake/dao"
	"testlake/model"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}

	dao.Connect()

	code := m.Run()
	os.Exit(code)
}

func TestUserDao_Create(t *testing.T) {
	userDao := dao.NewUserDao()

	user := &model.User{
		Email:        "test@example.com",
		Username:     "testuser",
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}

	err := userDao.Create(user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)

	defer userDao.Delete(user.ID)
}

func TestUserDao_GetByID(t *testing.T) {
	userDao := dao.NewUserDao()

	user := &model.User{
		Email:        "test2@example.com",
		Username:     "testuser2",
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}

	err := userDao.Create(user)
	assert.NoError(t, err)

	foundUser, err := userDao.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Username, foundUser.Username)

	defer userDao.Delete(user.ID)
}

func TestUserDao_GetByEmail(t *testing.T) {
	userDao := dao.NewUserDao()

	user := &model.User{
		Email:        "test3@example.com",
		Username:     "testuser3",
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}

	err := userDao.Create(user)
	assert.NoError(t, err)

	foundUser, err := userDao.GetByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)

	defer userDao.Delete(user.ID)
}

func TestUserDao_EmailExists(t *testing.T) {
	userDao := dao.NewUserDao()

	user := &model.User{
		Email:        "test4@example.com",
		Username:     "testuser4",
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}

	err := userDao.Create(user)
	assert.NoError(t, err)

	exists, err := userDao.EmailExists(user.Email)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = userDao.EmailExists("nonexistent@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)

	defer userDao.Delete(user.ID)
}

func TestUserDao_Update(t *testing.T) {
	userDao := dao.NewUserDao()

	user := &model.User{
		Email:        "test5@example.com",
		Username:     "testuser5",
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}

	err := userDao.Create(user)
	assert.NoError(t, err)

	firstName := "John"
	user.FirstName = &firstName

	err = userDao.Update(user)
	assert.NoError(t, err)

	updatedUser, err := userDao.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, firstName, *updatedUser.FirstName)

	defer userDao.Delete(user.ID)
}
