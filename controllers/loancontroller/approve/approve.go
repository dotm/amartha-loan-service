package approveloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ApproveLoanInput struct {
	LoanID                             string    `json:"loan_id" binding:"required"`
	FieldOfficerUserID                 string    `json:"field_officer_user_id" binding:"required"`
	VisitProofBeforeApprovalPictureUrl string    `json:"visit_proof_before_approval_picture_url" binding:"required"`
	TimeApproved                       time.Time `json:"time_approved" binding:"required"` //accepts ISO time format (2024-07-11T23:27:58.687Z)
}

func ApproveLoanHandler(c *gin.Context) {
	// Validate input
	var input ApproveLoanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get required data
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

	//Business logic
	updatedLoan, err := ApproveLoanBusinessLogic(input, user, loan, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Persist data
	models.DB.Updates(&updatedLoan)

	c.JSON(http.StatusOK, gin.H{"data": updatedLoan})
}

// Should be covered with unit test
func ApproveLoanBusinessLogic(input ApproveLoanInput, user models.User, loan models.Loan, now time.Time) (updatedLoan models.Loan, err error) {
	//Validate business logic
	if user.Role != constants.UserRoleFieldOfficer {
		err = fmt.Errorf("user should be %s", constants.UserRoleFieldOfficer)
		return
	}
	if loan.Status != constants.LoanStatusProposed {
		err = fmt.Errorf("loan status should be %s", constants.LoanStatusProposed)
		return
	}

	//Update data
	updatedLoan = loan
	updatedLoan.Status = constants.LoanStatusApproved
	updatedLoan.VisitProofBeforeApprovalPictureUrl = &input.VisitProofBeforeApprovalPictureUrl
	updatedLoan.TimeApproved = &input.TimeApproved
	return
}
