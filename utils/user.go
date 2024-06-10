package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetActiveUserId(c *gin.Context) (int64, error){
	userId, exists := c.Get("user_id")

	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return 0, fmt.Errorf("error occurred ")
	}
	
	convertedUserId, ok := userId.(int64)	

	if !ok{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Encountered an issue"})
		return 0, fmt.Errorf("error occurred ")
	}

	return convertedUserId, nil
}