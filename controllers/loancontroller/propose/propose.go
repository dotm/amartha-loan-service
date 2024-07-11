package proposeloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProposeLoanInput struct {
	BorrowerUserID       string  `json:"borrower_user_id" binding:"required"`
	PrincipalAmount      int64   `json:"principal_amount" binding:"required"`
	TenorInMonths        int64   `json:"tenor_in_months" binding:"required"`
	BorrowerRatePerMonth float64 `json:"borrower_rate_per_month" binding:"required"`
	InvestorRatePerMonth float64 `json:"investor_rate_per_month" binding:"required"`
	AgreementLetterUrl   string  `json:"agreement_letter_url" binding:"required"`
}

func ProposeLoanHandler(c *gin.Context) {
	// Validate input
	var input ProposeLoanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get required data
	var user models.User
	if err := models.DB.Where("id = ?", input.BorrowerUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	//Business logic
	loan, err := ProposeLoanBusinessLogic(input, user, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Persist data
	models.DB.Create(&loan)

	c.JSON(http.StatusOK, gin.H{"data": loan})
}

// Should be covered with unit test
func ProposeLoanBusinessLogic(input ProposeLoanInput, user models.User, now time.Time) (loan models.Loan, err error) {
	//Validate business logic
	if user.Role != constants.UserRoleBorrower {
		err = fmt.Errorf("user should be %s", constants.UserRoleBorrower)
		return
	}

	//Create data
	loan = models.Loan{
		ID:                   uuid.NewString(),
		BorrowerUserID:       input.BorrowerUserID,
		Status:               constants.LoanStatusProposed,
		PrincipalAmount:      input.PrincipalAmount,
		TenorInMonths:        input.TenorInMonths,
		BorrowerRatePerMonth: input.BorrowerRatePerMonth,
		InvestorRatePerMonth: input.InvestorRatePerMonth,
		AgreementLetterUrl:   input.AgreementLetterUrl,
		TimeProposed:         now,
	}
	return
}
