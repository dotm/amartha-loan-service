package approveloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/helpers/authorizationhelper"
	"amartha/loan-service/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ApproveLoanInput struct {
	LoanID       string    `json:"loan_id" binding:"required"`
	TimeApproved time.Time `json:"time_approved" binding:"required"` //accepts ISO time format (2024-07-11T23:27:58.687Z)

	VisitProofBeforeApprovalPictureUrl string `json:"visit_proof_before_approval_picture_url" binding:"required"`
}

func ApproveLoanHandler(c *gin.Context) {
	// Validate input
	var input ApproveLoanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Clean input
	// remove query params from S3 presigned URL
	input.VisitProofBeforeApprovalPictureUrl = strings.Split(input.VisitProofBeforeApprovalPictureUrl, "?")[0]

	// Get required data (can be moved to repository layer)
	fieldOfficerUserID, err := authorizationhelper.GetUserIdFromGinContextHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := models.DB.Where("id = ?", fieldOfficerUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	var loan models.Loan
	if err := models.DB.Where("id = ?", input.LoanID).First(&loan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan not found"})
		return
	}

	// Business Logic
	updatedLoan, err := ApproveLoanBusinessLogic(input, user, loan, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Persist data (can be moved to repository layer)
	// Use transaction for multiple update/create operations
	if err := models.DB.Updates(&updatedLoan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
