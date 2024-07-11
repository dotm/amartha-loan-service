package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=098098 dbname=amartha_loan_service port=5432 sslmode=disable TimeZone=Asia/Bangkok"

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
