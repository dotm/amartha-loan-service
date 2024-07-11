package approveloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var userID = "field_officer_1"
var principalAmount = int64(1000000)
var tenorInMonths = int64(12)
var borrowerRatePerMonth = 0.01
var investorRatePerMonth = 0.008
var agreementLetterUrl = "https://example.com"
var timeProposed = time.Now()

var loanID = "loan_1"
var borrowerUserID = "borrower_1"
var visitProofBeforeApprovalPictureUrl = "https://example.com/visit-proof"
var timeApproved = time.Now()

var input = ApproveLoanInput{
	LoanID:                             loanID,
	FieldOfficerUserID:                 userID,
	VisitProofBeforeApprovalPictureUrl: visitProofBeforeApprovalPictureUrl,
	TimeApproved:                       timeApproved,
}

// TestApproveLoanSuccessful is the happy path of propose loan
func TestApproveLoanSuccessful(t *testing.T) {
	//given
	proposedLoan := models.Loan{
		ID:                   loanID,
		BorrowerUserID:       borrowerUserID,
		Status:               constants.LoanStatusProposed,
		PrincipalAmount:      principalAmount,
		TenorInMonths:        tenorInMonths,
		BorrowerRatePerMonth: borrowerRatePerMonth,
		InvestorRatePerMonth: investorRatePerMonth,
		AgreementLetterUrl:   agreementLetterUrl,
		TimeProposed:         timeProposed,
	}
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleFieldOfficer,
		Name: "User 1",
	}

	//when
	loan, err := ApproveLoanBusinessLogic(input, user, proposedLoan, timeApproved)

	//then
	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, loan.Status, constants.LoanStatusApproved)
	assert.Equal(t, *loan.VisitProofBeforeApprovalPictureUrl, visitProofBeforeApprovalPictureUrl)
	assert.Equal(t, *loan.TimeApproved, timeApproved)
}

// TestApproveFailedIncorrectUserRole for the case when user is not field officer
func TestApproveFailedIncorrectUserRole(t *testing.T) {
	//given
	proposedLoan := models.Loan{
		ID:                   loanID,
		BorrowerUserID:       borrowerUserID,
		Status:               constants.LoanStatusProposed,
		PrincipalAmount:      principalAmount,
		TenorInMonths:        tenorInMonths,
		BorrowerRatePerMonth: borrowerRatePerMonth,
		InvestorRatePerMonth: investorRatePerMonth,
		AgreementLetterUrl:   agreementLetterUrl,
		TimeProposed:         timeProposed,
	}
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleBorrower,
		Name: "User 1",
	}

	//when
	_, err := ApproveLoanBusinessLogic(input, user, proposedLoan, timeApproved)

	//then
	assert.Error(t, err)
}

// TestApproveFailedIncorrectLoanStatus for the case when loan status is not proposed
func TestApproveFailedIncorrectLoanStatus(t *testing.T) {
	//given
	proposedLoan := models.Loan{
		ID:                   loanID,
		BorrowerUserID:       borrowerUserID,
		Status:               constants.LoanStatusApproved,
		PrincipalAmount:      principalAmount,
		TenorInMonths:        tenorInMonths,
		BorrowerRatePerMonth: borrowerRatePerMonth,
		InvestorRatePerMonth: investorRatePerMonth,
		AgreementLetterUrl:   agreementLetterUrl,
		TimeProposed:         timeProposed,
	}
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleFieldOfficer,
		Name: "User 1",
	}

	//when
	_, err := ApproveLoanBusinessLogic(input, user, proposedLoan, timeApproved)

	//then
	assert.Error(t, err)
}
