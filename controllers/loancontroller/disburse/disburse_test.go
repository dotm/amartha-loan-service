package disburseloancontroller

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
var timeProposed = time.Now().Add(-time.Minute)

var loanID = "loan_1"
var borrowerUserID = "borrower_1"
var signedLoanAgreementLetterUrl = "https://example.com/signed-loan"
var timeDisbursed = time.Now()

var input = DisburseLoanInput{
	LoanID:                       loanID,
	FieldOfficerUserID:           userID,
	SignedLoanAgreementLetterUrl: signedLoanAgreementLetterUrl,
	TimeDisbursed:                timeDisbursed,
}

// TestDisburseLoanSuccessful is the happy path of propose loan
func TestDisburseLoanSuccessful(t *testing.T) {
	//given
	investedLoan := models.Loan{
		ID:                   loanID,
		BorrowerUserID:       borrowerUserID,
		Status:               constants.LoanStatusInvested,
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
	loan, err := DisburseLoanBusinessLogic(input, user, investedLoan, timeDisbursed)

	//then
	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, constants.LoanStatusDisbursed, loan.Status)
	assert.Equal(t, signedLoanAgreementLetterUrl, *loan.SignedLoanAgreementLetterUrl)
	assert.Equal(t, timeDisbursed, *loan.TimeDisbursed)
}

// TestDisburseFailedIncorrectUserRole for the case when user is not field officer
func TestDisburseFailedIncorrectUserRole(t *testing.T) {
	//given
	investedLoan := models.Loan{
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
	_, err := DisburseLoanBusinessLogic(input, user, investedLoan, timeDisbursed)

	//then
	assert.Error(t, err)
}

// TestDisburseFailedIncorrectLoanStatus for the case when loan status is not proposed
func TestDisburseFailedIncorrectLoanStatus(t *testing.T) {
	//given
	investedLoan := models.Loan{
		ID:                   loanID,
		BorrowerUserID:       borrowerUserID,
		Status:               constants.LoanStatusDisbursed,
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
	_, err := DisburseLoanBusinessLogic(input, user, investedLoan, timeDisbursed)

	//then
	assert.Error(t, err)
}
