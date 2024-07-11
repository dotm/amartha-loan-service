package models

import "time"

type Loan struct {
	ID                   string     `json:"id" gorm:"primary_key"`
	BorrowerUserID       string     `json:"borrower_user_id"`
	Status               string     `json:"status"`
	PrincipalAmount      int64      `json:"principal_amount"`
	TenorInMonths        int64      `json:"tenor_in_months"`
	BorrowerRatePerMonth float64    `json:"borrower_rate_per_month"`
	InvestorRatePerMonth float64    `json:"investor_rate_per_month"`
	AgreementLetterUrl   string     `json:"agreement_letter_url"`
	TimeProposed         time.Time  `json:"time_proposed"`
	TimeApproved         *time.Time `json:"time_approved"`
	TimeInvested         *time.Time `json:"time_invested"` //when the whole principal amount is invested
	TimeDisbursed        *time.Time `json:"time_disbursed"`

	VisitProofBeforeApprovalPictureUrl *string `json:"visit_proof_before_approval_picture_url"`
	SignedLoanAgreementLetterUrl       *string `json:"signed_loan_agreement_letter_url"`

	Investor []LoanInvestedByInvestor `json:"investor" gorm:"foreignKey:LoanID"`
}

type LoanInvestedByInvestor struct {
	ID             string    `json:"id" gorm:"primary_key"`
	InvestorUserID string    `json:"investor_user_id"`
	LoanID         string    `json:"loan_id"`
	Amount         int64     `json:"amount"`
	TimeInvested   time.Time `json:"time_invested"` //when investor chip into the loan
}
