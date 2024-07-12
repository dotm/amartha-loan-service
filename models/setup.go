package models

import (
	"amartha/loan-service/helpers/envhelper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := envhelper.GetEnvVar("POSTGRESQL_DSN")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&User{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&Loan{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&LoanInvestedByInvestor{})
	if err != nil {
		return
	}

	DB = database
}
