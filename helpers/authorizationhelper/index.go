package authorizationhelper

import (
	"amartha/loan-service/helpers/jwthelper"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromGinContextHeader(c *gin.Context) (userID string, err error) {
	authorizationHeaders := c.Request.Header["Authorization"]
	if len(authorizationHeaders) == 0 {
		return "", fmt.Errorf("empty authorization header")
	}
	authorizationHeader := authorizationHeaders[0]
	token, err := jwthelper.ParseFromAuthorizationHeader(authorizationHeader)
	if err != nil {
		return "", err
	}

	return jwthelper.GetUserId(token), nil
}
