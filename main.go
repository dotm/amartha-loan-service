package main

import (
	approveloancontroller "amartha/loan-service/controllers/loancontroller/approve"
	proposeloancontroller "amartha/loan-service/controllers/loancontroller/propose"
	"amartha/loan-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	r.POST("/v1/loan.propose", proposeloancontroller.ProposeLoanHandler)
	r.POST("/v1/loan.approve", approveloancontroller.ApproveLoanHandler)

	r.Run()
}
