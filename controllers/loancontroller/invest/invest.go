package investloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/helpers/authorizationhelper"
	"amartha/loan-service/helpers/emailhelper"
	"amartha/loan-service/helpers/envhelper"
	"amartha/loan-service/helpers/pdfhelper"
	"amartha/loan-service/helpers/s3helper"
	"amartha/loan-service/models"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvestLoanInput struct {
	LoanID string `json:"loan_id" binding:"required"`
	Amount int64  `json:"amount" binding:"required"`
}

func InvestLoanHandler(c *gin.Context) {
	// Validate input
	var input InvestLoanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get required data (can be moved to repository layer)
	investorUserID, err := authorizationhelper.GetUserIdFromGinContextHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := models.DB.Where("id = ?", investorUserID).First(&user).Error; err != nil {
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

	// Business Logic
	updatedLoan, newInvestment, err := InvestLoanBusinessLogic(input, user, loan, investments, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Interact with External APIs as side effect
	if updatedLoan.Status == constants.LoanStatusInvested {
		// Get all investor emails
		investorIDs := []string{newInvestment.InvestorUserID}
		for _, investment := range investments {
			investorIDs = append(investorIDs, investment.InvestorUserID)
		}
		var investors []models.User
		if err := models.DB.Where("id IN ?", investorIDs).Find(&investors).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		investorEmails := []string{}
		for _, investor := range investors {
			investorEmails = append(investorEmails, investor.Email)
		}

		s3Client, err := s3helper.CreateClientFromSession()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sesClient, err := emailhelper.CreateClientFromSession()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		presignedUrlForGet, err := GenerateAndUploadPDF(s3Client, input, loan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, investorEmail := range investorEmails {
			if err != nil {
				fmt.Println(err)
			}
			body := fmt.Sprintf(
				"A loan you have invested in has been fully invested and ready to be disbursed. This is the link to the unsigned agreement letter (expired in %s): %s",
				s3helper.DefaultExpireDurationForPresignedURLString,
				presignedUrlForGet,
			)
			err = emailhelper.SendEmail(sesClient, emailhelper.SendEmailArgs{
				Sender:        emailhelper.DefaultEmailSender,
				RecipientList: []*string{aws.String(investorEmail)},
				CcList:        []*string{},
				Subject:       "Loan Invested",
				HtmlBody:      body,
				TextBody:      body,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
				//we can also store errors in an array and continue with the next email instead of terminating early
			}
		}
	}

	// Persist data (can be moved to repository layer)
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
) (
	updatedLoan models.Loan,
	newInvestment models.LoanInvestedByInvestor,
	err error,
) {
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
	userIsExistingInvestor := false
	for i := 0; i < len(existingInvestments); i++ {
		existingInvestmentAmount += existingInvestments[i].Amount
		if existingInvestments[i].InvestorUserID == user.ID {
			userIsExistingInvestor = true
		}
	}
	if userIsExistingInvestor {
		//can be improved later to allow logic branch for updating or adding new amount to a loan
		err = fmt.Errorf("you can only invest once in a loan")
		return
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

func GenerateAndUploadPDF(
	s3Client *s3.S3,
	input InvestLoanInput,
	loan models.Loan,
) (presignedUrlForGet string, err error) {
	pdfFile := pdfhelper.GenerateLoanAgreementLetter(pdfhelper.GenerateLoanAgreementLetterArgs{
		PrincipalAmount: loan.PrincipalAmount,
		TenorInMonths:   loan.TenorInMonths,
		BorrowerName:    "Borrower 1",
	})
	buf := new(bytes.Buffer)
	err = pdfFile.Output(buf)
	if err != nil {
		err = fmt.Errorf("error writing pdf to pipe: %v", err)
		return "", err
	}
	keyName := s3helper.GetAgreementLetterKeyName(input.LoanID, false, "pdf")
	err = s3helper.UploadFile(
		s3Client,
		buf,
		envhelper.GetEnvVar("S3_BUCKET_NAME"),
		keyName,
	)
	if err != nil {
		return "", err
	}

	presignedUrl, err := s3helper.GeneratePresignedURLForGetObject(
		s3Client,
		envhelper.GetEnvVar("S3_BUCKET_NAME"),
		keyName,
	)
	if err != nil {
		return "", err
	}

	return presignedUrl, nil
}
