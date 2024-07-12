package main

import (
	approveloancontroller "amartha/loan-service/controllers/loancontroller/approve"
	disburseloancontroller "amartha/loan-service/controllers/loancontroller/disburse"
	investloancontroller "amartha/loan-service/controllers/loancontroller/invest"
	proposeloancontroller "amartha/loan-service/controllers/loancontroller/propose"
	uploadsignedagreementcontroller "amartha/loan-service/controllers/objectstoragecontroller/getpresignedurl/uploadsignedagreement"
	uploadvisitproofcontroller "amartha/loan-service/controllers/objectstoragecontroller/getpresignedurl/uploadvisitproof"
	usersignincontroller "amartha/loan-service/controllers/usercontroller/signin"
	usersignupcontroller "amartha/loan-service/controllers/usercontroller/signup"
	"amartha/loan-service/helpers/envhelper"
	"amartha/loan-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	envhelper.SetLocalEnvVar()
	r := gin.Default()
	models.ConnectDatabase()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	//The endpoints below use slack convention, but I can use proper REST if it's the codebase convention
	r.POST("/v1/user.signup", usersignupcontroller.UserSignUpHandler)
	r.POST("/v1/user.signin", usersignincontroller.UserSignInHandler)
	r.POST("/v1/loan.propose", proposeloancontroller.ProposeLoanHandler)
	r.POST("/v1/loan.approve", approveloancontroller.ApproveLoanHandler)
	r.POST("/v1/loan.invest", investloancontroller.InvestLoanHandler)
	r.POST("/v1/loan.disburse", disburseloancontroller.DisburseLoanHandler)
	//using POST because the endpoint already specify action (but I'm open to using GET if the codebase convention is to use proper REST)
	r.POST("/v1/object_storage.get_presigned_url.upload_visit_proof", uploadvisitproofcontroller.UploadVisitProofHandler)
	r.POST("/v1/object_storage.get_presigned_url.upload_signed_agreement", uploadsignedagreementcontroller.UploadSignedAgreementHandler)

	r.Run()
}
