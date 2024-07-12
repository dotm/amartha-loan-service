package usersignincontroller

import (
	"amartha/loan-service/helpers/jwthelper"
	"amartha/loan-service/helpers/passwordhelper"
	"amartha/loan-service/models"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserSignInInput struct {
	Email       string `json:"email" binding:"required"`
	RawPassword string `json:"raw_password" binding:"required"`
}

func UserSignInHandler(c *gin.Context) {
	// Validate input
	var input UserSignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user is registered
	var existingUser models.User
	err := models.DB.Where("email = ?", input.Email).First(&existingUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user email has not been registered"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed checking if user is registered: %v", err)})
		return
	}

	// Business Logic
	//check password is correct
	passwordCorrect := passwordhelper.MatchPasswordToHash(input.RawPassword, existingUser.HashedPassword)
	if !passwordCorrect {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user password incorrect"})
		return
	}
	jwtExpiration := time.Now().Add(time.Hour * 24 * 365) //365 days
	signedJwtToken, err := jwthelper.BuildAndSign(
		jwthelper.BuildCustomClaims(existingUser.ID),
		jwtExpiration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signedJwtToken})
}
