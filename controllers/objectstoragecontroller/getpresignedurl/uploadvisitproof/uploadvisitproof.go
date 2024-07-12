package uploadvisitproofcontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/helpers/authorizationhelper"
	"amartha/loan-service/helpers/envhelper"
	"amartha/loan-service/helpers/s3helper"
	"amartha/loan-service/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadVisitProofInput struct {
	LoanID string `json:"loan_id" binding:"required"`
}

func UploadVisitProofHandler(c *gin.Context) {
	// Validate input
	var input UploadVisitProofInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	if user.Role != constants.UserRoleFieldOfficer {
		err = fmt.Errorf("user should be %s", constants.UserRoleFieldOfficer)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if loan.Status != constants.LoanStatusProposed {
		err = fmt.Errorf("loan status should be %s", constants.LoanStatusProposed)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Presigned URL
	s3Client, err := s3helper.CreateClientFromSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	presignedURLForPutObject, err := s3helper.GeneratePresignedURLForPutObject(
		s3Client,
		envhelper.GetEnvVar("S3_BUCKET_NAME"),
		fmt.Sprintf("visit-proof/%s", input.LoanID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": presignedURLForPutObject})
}
