package usersignupcontroller

import (
	"amartha/loan-service/constants"
	"amartha/loan-service/helpers/passwordhelper"
	"amartha/loan-service/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSignUpInput struct {
	Role        string `json:"role" binding:"required"`
	Email       string `json:"email" binding:"required"`
	RawPassword string `json:"raw_password" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

func UserSignUpHandler(c *gin.Context) {
	// Validate input
	var input UserSignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user not registered
	var existingUser models.User
	err := models.DB.Where("email = ?", input.Email).First(&existingUser).Error
	if err == nil || existingUser.Email == input.Email {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user email already registered"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed checking if user is registered: %v", err)})
		return
	}

	// Business Logic
	if input.Role != constants.UserRoleBorrower &&
		input.Role != constants.UserRoleInvestor &&
		input.Role != constants.UserRoleFieldOfficer {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user role invalid"})
		return
	}
	hashedPassword, err := passwordhelper.Hash(input.RawPassword)
	if err != nil {
		err = fmt.Errorf("error hashing password: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser := models.User{
		ID:             uuid.NewString(),
		Role:           input.Role,
		Email:          input.Email,
		Name:           input.Name,
		HashedPassword: hashedPassword,
	}

	// Persist data (can be moved to repository layer)
	models.DB.Create(&newUser)

	c.JSON(http.StatusOK, gin.H{"data": newUser})
}
