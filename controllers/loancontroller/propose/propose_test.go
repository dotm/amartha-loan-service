package proposeloancontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var userID = "borrower_1"
var principalAmount = int64(1000000)
var tenorInMonths = int64(12)
var borrowerRatePerMonth = 0.01
var investorRatePerMonth = 0.008
var agreementLetterUrl = "https://example.com"
var timeProposed = time.Now()

var input = ProposeLoanInput{
	PrincipalAmount:      principalAmount,
	TenorInMonths:        tenorInMonths,
	BorrowerRatePerMonth: borrowerRatePerMonth,
	InvestorRatePerMonth: investorRatePerMonth,
	AgreementLetterUrl:   agreementLetterUrl,
}

// TestProposeLoanSuccessful is the happy path of propose loan
func TestProposeLoanSuccessful(t *testing.T) {
	//given
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleBorrower,
		Name: "Borrower 1",
	}

	//when
	loan, err := ProposeLoanBusinessLogic(input, user, timeProposed)

	//then
	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, userID, loan.BorrowerUserID)
	assert.Equal(t, constants.LoanStatusProposed, loan.Status)
	assert.Equal(t, principalAmount, loan.PrincipalAmount)
	assert.Equal(t, tenorInMonths, loan.TenorInMonths)
	assert.Equal(t, borrowerRatePerMonth, loan.BorrowerRatePerMonth)
	assert.Equal(t, investorRatePerMonth, loan.InvestorRatePerMonth)
	assert.Equal(t, agreementLetterUrl, loan.AgreementLetterUrl)
	assert.Equal(t, timeProposed, loan.TimeProposed)
}

// TestProposeFailedIncorrectUserRole for the case when user is not borrower
func TestProposeFailedIncorrectUserRole(t *testing.T) {
	//given
	user := models.User{
		ID:   userID,
		Role: constants.UserRoleInvestor,
		Name: "Investor 1",
	}

	//when
	_, err := ProposeLoanBusinessLogic(input, user, timeProposed)

	//then
	assert.Error(t, err)
}
