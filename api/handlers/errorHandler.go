package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator/v2"
)

func HandleError(err error, c *gin.Context, gValid galidator.Validator) interface{}{
	if c.Request.ContentLength == 0{
		return "provide body"
	}

	if e, ok := err.(*json.UnmarshalTypeError); ok{
		if e.Field == ""{
			return "provide json body"
		}
		msg := fmt.Sprintf("invalid type for field %s. Expectec value of type: %s", e.Field, e.Type)
		return msg
	}
	return gValid.DecryptErrors(err)

}