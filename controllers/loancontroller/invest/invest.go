package investloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/models"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvestLoanInput struct {
	LoanID         string `json:"loan_id" binding:"required"`
	InvestorUserID string `json:"investor_user_id" binding:"required"`
	Amount         int64  `json:"amount" binding:"required"`
}

func InvestLoanHandler(c *gin.Context) {
	// Validate input
	var input InvestLoanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get required data
	var user models.User
	if err := models.DB.Where("id = ?", input.InvestorUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	var loan models.Loan
	if err := models.DB.Where("id = ?", input.LoanID).First(&loan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan not found"})
		return
	}
	var investments []models.LoanInvestedByInvestor
	if err := models.DB.Where("loan_id = ?", input.LoanID).Find(&investments).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Business logic
	updatedLoan, newInvestment, err := InvestLoanBusinessLogic(input, user, loan, investments, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Persist data
	models.DB.Updates(&updatedLoan)
	models.DB.Create(&newInvestment)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"loan":           updatedLoan,
			"new_investment": newInvestment,
		},
	})
}

// Should be covered with unit test
func InvestLoanBusinessLogic(
	input InvestLoanInput,
	user models.User,
	loan models.Loan,
	existingInvestments []models.LoanInvestedByInvestor,
	now time.Time,
) (updatedLoan models.Loan, newInvestment models.LoanInvestedByInvestor, err error) {
	//Validate business logic
	if user.Role != constants.UserRoleInvestor {
		err = fmt.Errorf("user should be %s", constants.UserRoleInvestor)
		return
	}
	if loan.Status != constants.LoanStatusApproved {
		err = fmt.Errorf("loan status should be %s", constants.LoanStatusApproved)
		return
	}
	existingInvestmentAmount := int64(0)
	for i := 0; i < len(existingInvestments); i++ {
		existingInvestmentAmount += existingInvestments[i].Amount
	}
	if existingInvestmentAmount+input.Amount > loan.PrincipalAmount {
		err = fmt.Errorf("maximum amount available for investment is %d", loan.PrincipalAmount-existingInvestmentAmount)
		return
	}

	//Update data
	updatedLoan = loan
	principalIsFulfilled := loan.PrincipalAmount == existingInvestmentAmount+input.Amount
	newInvestment = models.LoanInvestedByInvestor{
		ID:             uuid.NewString(),
		InvestorUserID: user.ID,
		LoanID:         loan.ID,
		Amount:         input.Amount,
		TimeInvested:   now,
	}
	if principalIsFulfilled {
		updatedLoan.Status = constants.LoanStatusInvested
		updatedLoan.TimeInvested = &now
	}
	return
}
