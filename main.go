package main

import (
	loancontroller "amartha/loan-service/controllers/loancontroller"
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
	r.POST("/v1/loan.propose", loancontroller.ProposeLoanHandler)

	r.Run()
}
