package disburseloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DisburseLoanInput struct {
	LoanID                       string    `json:"loan_id" binding:"required"`
	FieldOfficerUserID           string    `json:"field_officer_user_id" binding:"required"`
	SignedLoanAgreementLetterUrl string    `json:"signed_loan_agreement_letter_url" binding:"required"`
	TimeDisbursed                time.Time `json:"time_disbursed" binding:"required"` //accepts ISO time format (2024-07-11T23:27:58.687Z)
}

func DisburseLoanHandler(c *gin.Context) {
	// Validate input
	var input DisburseLoanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get required data (can be moved to repository layer)
	var user models.User
	if err := models.DB.Where("id = ?", input.FieldOfficerUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	var loan models.Loan
	if err := models.DB.Where("id = ?", input.LoanID).First(&loan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan not found"})
		return
	}

	// Business Logic
	updatedLoan, err := DisburseLoanBusinessLogic(input, user, loan, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Persist data (can be moved to repository layer)
	models.DB.Updates(&updatedLoan)

	c.JSON(http.StatusOK, gin.H{"data": updatedLoan})
}

// Should be covered with unit test
func DisburseLoanBusinessLogic(input DisburseLoanInput, user models.User, loan models.Loan, now time.Time) (updatedLoan models.Loan, err error) {
	//Validate business logic
	if user.Role != constants.UserRoleFieldOfficer {
		err = fmt.Errorf("user should be %s", constants.UserRoleFieldOfficer)
		return
	}
	if loan.Status != constants.LoanStatusInvested {
		err = fmt.Errorf("loan status should be %s", constants.LoanStatusInvested)
		return
	}

	//Update data
	updatedLoan = loan
	updatedLoan.Status = constants.LoanStatusDisbursed
	updatedLoan.SignedLoanAgreementLetterUrl = &input.SignedLoanAgreementLetterUrl
	updatedLoan.TimeDisbursed = &input.TimeDisbursed
	return
}
