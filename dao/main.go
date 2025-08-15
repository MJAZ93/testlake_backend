package dao

import (
	"fmt"
	"log"
	"os"

	"testlake/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, username, password, databaseName, port)

	var err error
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = Database.AutoMigrate(
		&model.User{},
		&model.Organization{},
		&model.OrganizationMember{},
		&model.OrganizationInvitation{},
		&model.Project{},
		&model.Team{},
		&model.TeamMember{},
		&model.ProjectAccess{},
		&model.Environment{},
		&model.Feature{},
		&model.FeatureEnvironmentStatus{},
		&model.FeatureErrorLog{},
		&model.ErrorImage{},
		&model.DataSchema{},
		&model.FeatureSchema{},
		&model.SchemaField{},
		&model.TestData{},
		&model.TestDataRequest{},
		&model.EmailVerificationToken{},
		&model.PaymentMethod{},
		&model.Subscription{},
	)
	if err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}
}
