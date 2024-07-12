package investloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var userID = "investor_1"
var principalAmount = int64(1000000)
var tenorInMonths = int64(12)
var borrowerRatePerMonth = 0.01
var investorRatePerMonth = 0.008
var agreementLetterUrl = "https://example.com"
var timeProposed = time.Now().Add(-time.Minute * 2)
var timeApproved = time.Now().Add(-time.Minute)
var visitProofBeforeApprovalPictureUrl = "https://example.com/visit-proof"

var loanID = "loan_1"
var borrowerUserID = "borrower_1"
var timeInvested = time.Now()

var existingInvestments = []models.LoanInvestedByInvestor{
	{
		ID:             "previous-investment-1",
		InvestorUserID: "previous-investor-1",
		LoanID:         loanID,
		Amount:         int64(200000),
		TimeInvested:   time.Now().Add(-time.Hour),
	},
	{
		ID:             "previous-investment-2",
		InvestorUserID: "previous-investor-2",
		LoanID:         loanID,
		Amount:         int64(300000),
		TimeInvested:   time.Now().Add(-time.Hour),
	},
}

// TestInvestLoanSuccessfulFullyInvested for the case when investment has reach the principal amount
func TestInvestLoanSuccessfulFullyInvested(t *testing.T) {
	//given
	input := InvestLoanInput{
		LoanID: loanID,
		Amount: int64(500000),
	}
	approvedLoan := models.Loan{
		ID:                                 loanID,
		BorrowerUserID:                     borrowerUserID,
		Status:                             constants.LoanStatusApproved,
		PrincipalAmount:                    principalAmount,
		TenorInMonths:                      tenorInMonths,
		BorrowerRatePerMonth:               borrowerRatePerMonth,
		InvestorRatePerMonth:               investorRatePerMonth,
		AgreementLetterUrl:                 agreementLetterUrl,
		VisitProofBeforeApprovalPictureUrl: &visitProofBeforeApprovalPictureUrl,
		TimeProposed:                       timeProposed,
		TimeApproved:                       &timeApproved,
	}
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleInvestor,
		Name: "User 1",
	}

	//when
	loan, newInvestment, err := InvestLoanBusinessLogic(input, user, approvedLoan, existingInvestments, timeInvested)

	//then
	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, input.LoanID, loan.ID)
	assert.Equal(t, constants.LoanStatusInvested, loan.Status)
	assert.Equal(t, timeInvested, *loan.TimeInvested)
	assert.Equal(t, input.Amount, newInvestment.Amount)
	assert.Equal(t, user.ID, newInvestment.InvestorUserID)
	assert.Equal(t, input.LoanID, newInvestment.LoanID)
}

// TestInvestLoanSuccessfulButNotFullyInvested for the case when investment still hasn't reach the principal amount
func TestInvestLoanSuccessfulButNotFullyInvested(t *testing.T) {
	//given
	input := InvestLoanInput{
		LoanID: loanID,
		Amount: int64(100000),
	}
	approvedLoan := models.Loan{
		ID:                                 loanID,
		BorrowerUserID:                     borrowerUserID,
		Status:                             constants.LoanStatusApproved,
		PrincipalAmount:                    principalAmount,
		TenorInMonths:                      tenorInMonths,
		BorrowerRatePerMonth:               borrowerRatePerMonth,
		InvestorRatePerMonth:               investorRatePerMonth,
		AgreementLetterUrl:                 agreementLetterUrl,
		VisitProofBeforeApprovalPictureUrl: &visitProofBeforeApprovalPictureUrl,
		TimeProposed:                       timeProposed,
		TimeApproved:                       &timeApproved,
	}
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleInvestor,
		Name: "User 1",
	}

	//when
	loan, newInvestment, err := InvestLoanBusinessLogic(input, user, approvedLoan, existingInvestments, timeInvested)

	//then
	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, input.LoanID, loan.ID)
	assert.Equal(t, constants.LoanStatusApproved, loan.Status)
	assert.Equal(t, input.Amount, newInvestment.Amount)
	assert.Equal(t, user.ID, newInvestment.InvestorUserID)
	assert.Equal(t, input.LoanID, newInvestment.LoanID)
}

// TestInvestFailedOverSubscribed for the case when invested amount exceeds principal amount
func TestInvestFailedOverSubscribed(t *testing.T) {
	//given
	input := InvestLoanInput{
		LoanID: loanID,
		Amount: int64(900000),
	}
	approvedLoan := models.Loan{
		ID:                                 loanID,
		BorrowerUserID:                     borrowerUserID,
		Status:                             constants.LoanStatusApproved,
		PrincipalAmount:                    principalAmount,
		TenorInMonths:                      tenorInMonths,
		BorrowerRatePerMonth:               borrowerRatePerMonth,
		InvestorRatePerMonth:               investorRatePerMonth,
		AgreementLetterUrl:                 agreementLetterUrl,
		VisitProofBeforeApprovalPictureUrl: &visitProofBeforeApprovalPictureUrl,
		TimeProposed:                       timeProposed,
		TimeApproved:                       &timeApproved,
	}
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleInvestor,
		Name: "User 1",
	}

	//when
	_, _, err := InvestLoanBusinessLogic(input, user, approvedLoan, existingInvestments, timeInvested)

	//then
	assert.Error(t, err)
}

// TestInvestFailedUserAlreadyInvest for the case when user already invest before
func TestInvestFailedUserAlreadyInvest(t *testing.T) {
	//given
	input := InvestLoanInput{
		LoanID: loanID,
		Amount: int64(100000),
	}
	approvedLoan := models.Loan{
		ID:                                 loanID,
		BorrowerUserID:                     borrowerUserID,
		Status:                             constants.LoanStatusApproved,
		PrincipalAmount:                    principalAmount,
		TenorInMonths:                      tenorInMonths,
		BorrowerRatePerMonth:               borrowerRatePerMonth,
		InvestorRatePerMonth:               investorRatePerMonth,
		AgreementLetterUrl:                 agreementLetterUrl,
		VisitProofBeforeApprovalPictureUrl: &visitProofBeforeApprovalPictureUrl,
		TimeProposed:                       timeProposed,
		TimeApproved:                       &timeApproved,
	}
	user := models.User{
		ID:   "previous-investor-1",
		Role: constants.UserRoleInvestor,
		Name: "User 1",
	}

	//when
	_, _, err := InvestLoanBusinessLogic(input, user, approvedLoan, existingInvestments, timeInvested)

	//then
	assert.Error(t, err)
}

// TestInvestFailedIncorrectUserRole for the case when user is not investor
func TestInvestFailedIncorrectUserRole(t *testing.T) {
	//given
	input := InvestLoanInput{
		LoanID: loanID,
		Amount: int64(100000),
	}
	approvedLoan := models.Loan{
		ID:                                 loanID,
		BorrowerUserID:                     borrowerUserID,
		Status:                             constants.LoanStatusApproved,
		PrincipalAmount:                    principalAmount,
		TenorInMonths:                      tenorInMonths,
		BorrowerRatePerMonth:               borrowerRatePerMonth,
		InvestorRatePerMonth:               investorRatePerMonth,
		AgreementLetterUrl:                 agreementLetterUrl,
		VisitProofBeforeApprovalPictureUrl: &visitProofBeforeApprovalPictureUrl,
		TimeProposed:                       timeProposed,
		TimeApproved:                       &timeApproved,
	}
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleBorrower,
		Name: "User 1",
	}

	//when
	_, _, err := InvestLoanBusinessLogic(input, user, approvedLoan, existingInvestments, timeInvested)

	//then
	assert.Error(t, err)
}

// TestInvestFailedIncorrectLoanStatus for the case when loan status is not approved
func TestInvestFailedIncorrectLoanStatus(t *testing.T) {
	//given
	input := InvestLoanInput{
		LoanID: loanID,
		Amount: int64(100000),
	}
	approvedLoan := models.Loan{
		ID:                                 loanID,
		BorrowerUserID:                     borrowerUserID,
		Status:                             constants.LoanStatusProposed,
		PrincipalAmount:                    principalAmount,
		TenorInMonths:                      tenorInMonths,
		BorrowerRatePerMonth:               borrowerRatePerMonth,
		InvestorRatePerMonth:               investorRatePerMonth,
		AgreementLetterUrl:                 agreementLetterUrl,
		VisitProofBeforeApprovalPictureUrl: &visitProofBeforeApprovalPictureUrl,
		TimeProposed:                       timeProposed,
		TimeApproved:                       &timeApproved,
	}
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleInvestor,
		Name: "User 1",
	}

	//when
	_, _, err := InvestLoanBusinessLogic(input, user, approvedLoan, existingInvestments, timeInvested)

	//then
	assert.Error(t, err)
}
