package models

type User struct {
	ID             string `json:"id" gorm:"primary_key"`
	Role           string `json:"role"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	Name           string `json:"name"`

	InvestIn []LoanInvestedByInvestor `json:"investor" gorm:"foreignKey:InvestorUserID"`
}
