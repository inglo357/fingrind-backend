package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token == ""{
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			ctx.Abort()
			return
		}

		tokenSplit := strings.Split(token, " ")

		if len(tokenSplit) != 2 || strings.ToLower(tokenSplit[0]) != "bearer"{
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			ctx.Abort()
			return
		}

		userId, err := TokenHandler.VerifyToken(string(tokenSplit[1]))

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userId)
	}
}