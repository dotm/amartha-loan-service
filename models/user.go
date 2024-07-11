package models

type User struct {
	ID   string `json:"id" gorm:"primary_key"`
	Role string `json:"role"`
	Name string `json:"name"`

	InvestIn []LoanInvestedByInvestor `json:"investor" gorm:"foreignKey:InvestorUserID"`
}
