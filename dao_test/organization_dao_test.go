package dao_test

import (
	"fmt"
	"testing"
	"time"

	"testlake/dao"
	"testlake/model"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestOrganizationDao(t *testing.T) {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// Connect to database
	dao.Connect()

	// Create test user first
	userDao := dao.NewUserDao()
	timestamp := time.Now().UnixNano()
	testUser := &model.User{
		ID:           uuid.New(),
		Email:        fmt.Sprintf("test%d@example.com", timestamp),
		Username:     fmt.Sprintf("testuser%d", timestamp),
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}
	err = userDao.Create(testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test organization creation
	orgDao := dao.NewOrganizationDao()
	testOrg := &model.Organization{
		Name:        fmt.Sprintf("Test Organization %d", timestamp),
		Slug:        fmt.Sprintf("test-org-%d", timestamp),
		Description: stringPtr("A test organization"),
		CreatedBy:   testUser.ID,
		PlanType:    model.PlanTypeStarter,
		Status:      model.OrganizationStatusActive,
	}

	err = orgDao.Create(testOrg)
	if err != nil {
		t.Fatalf("Failed to create organization: %v", err)
	}

	// Test get by ID
	foundOrg, err := orgDao.GetByID(testOrg.ID)
	if err != nil {
		t.Fatalf("Failed to get organization by ID: %v", err)
	}
	if foundOrg.Name != testOrg.Name {
		t.Errorf("Expected name %s, got %s", testOrg.Name, foundOrg.Name)
	}

	// Test get by slug
	foundOrgBySlug, err := orgDao.GetBySlug(testOrg.Slug)
	if err != nil {
		t.Fatalf("Failed to get organization by slug: %v", err)
	}
	if foundOrgBySlug.ID != testOrg.ID {
		t.Errorf("Expected ID %s, got %s", testOrg.ID, foundOrgBySlug.ID)
	}

	// Test update
	testOrg.Description = stringPtr("Updated description")
	err = orgDao.Update(testOrg)
	if err != nil {
		t.Fatalf("Failed to update organization: %v", err)
	}

	// Verify update
	updatedOrg, err := orgDao.GetByID(testOrg.ID)
	if err != nil {
		t.Fatalf("Failed to get updated organization: %v", err)
	}
	if *updatedOrg.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got %s", *updatedOrg.Description)
	}

	// Test delete
	err = orgDao.Delete(testOrg.ID)
	if err != nil {
		t.Fatalf("Failed to delete organization: %v", err)
	}

	// Verify deletion
	_, err = orgDao.GetByID(testOrg.ID)
	if err == nil {
		t.Error("Expected error when getting deleted organization, got none")
	}

	// Clean up test user
	userDao.Delete(testUser.ID)
}

func stringPtr(s string) *string {
	return &s
}