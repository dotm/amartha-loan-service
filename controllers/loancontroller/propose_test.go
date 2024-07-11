package loancontroller_test

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/controllers/loancontroller"
	"amartha/loan-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestProposeLoanSuccessful is the happy path of propose loan
func TestProposeLoanSuccessful(t *testing.T) {
	//given
	userId := "borrower_1"
	principalAmount := int64(1000000)
	tenorInMonths := int64(12)
	borrowerRatePerMonth := 0.01
	investorRatePerMonth := 0.008
	agreementLetterUrl := "https://example.com"
	timeProposed := time.Now()
	input := loancontroller.ProposeLoanInput{
		BorrowerUserID:       userId,
		PrincipalAmount:      principalAmount,
		TenorInMonths:        tenorInMonths,
		BorrowerRatePerMonth: borrowerRatePerMonth,
		InvestorRatePerMonth: investorRatePerMonth,
		AgreementLetterUrl:   agreementLetterUrl,
	}
	user := models.User{
		ID:   userId,
		Role: constants.UserRoleBorrower,
		Name: "Borrower 1",
	}

	//when
	loan, err := loancontroller.ProposeLoanBusinessLogic(input, user, timeProposed)

	//then
	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, loan.BorrowerUserID, userId)
	assert.Equal(t, loan.Status, constants.LoanStatusProposed)
	assert.Equal(t, loan.PrincipalAmount, principalAmount)
	assert.Equal(t, loan.TenorInMonths, tenorInMonths)
	assert.Equal(t, loan.BorrowerRatePerMonth, borrowerRatePerMonth)
	assert.Equal(t, loan.InvestorRatePerMonth, investorRatePerMonth)
	assert.Equal(t, loan.AgreementLetterUrl, agreementLetterUrl)
	assert.Equal(t, loan.TimeProposed, timeProposed)
}

// TestProposeFailedIncorrectRole for the case when user is not borrower
func TestProposeFailedIncorrectRole(t *testing.T) {
	//given
	userId := "investor_1"
	principalAmount := int64(1000000)
	tenorInMonths := int64(12)
	borrowerRatePerMonth := 0.01
	investorRatePerMonth := 0.008
	agreementLetterUrl := "https://example.com"
	timeProposed := time.Now()
	input := loancontroller.ProposeLoanInput{
		BorrowerUserID:       userId,
		PrincipalAmount:      principalAmount,
		TenorInMonths:        tenorInMonths,
		BorrowerRatePerMonth: borrowerRatePerMonth,
		InvestorRatePerMonth: investorRatePerMonth,
		AgreementLetterUrl:   agreementLetterUrl,
	}
	user := models.User{
		ID:   userId,
		Role: constants.UserRoleInvestor,
		Name: "Investor 1",
	}

	//when
	_, err := loancontroller.ProposeLoanBusinessLogic(input, user, timeProposed)

	//then
	assert.Error(t, err)
}
