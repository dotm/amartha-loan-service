package pdfhelper

import (
	"fmt"

	"github.com/go-pdf/fpdf"
)

type GenerateLoanAgreementLetterArgs struct {
	PrincipalAmount int64
	TenorInMonths   int64
	BorrowerName    string
}

func GenerateLoanAgreementLetter(args GenerateLoanAgreementLetterArgs) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Loan Agreement Letter", "", 1, "", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 10, fmt.Sprintf("Principal Amount: %d", args.PrincipalAmount), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Tenor: %d month(s)", args.TenorInMonths), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Terms and Conditions:", "", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "...", "", 1, "", false, 0, "")
	pdf.CellFormat(0, 30, fmt.Sprintf("Signed by %s:", args.BorrowerName), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 40, "_________________", "", 1, "", false, 0, "")
	return pdf
}
