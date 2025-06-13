package common

import (
	"github.com/twoonefour/alist-auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Error(c *gin.Context, err error) {
	utils.GetLogger(c).Error("[ERROR] %v\n", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}

func ErrorStr(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

func ErrorJson(c *gin.Context, err interface{}, code ...int) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		Code: func() int {
			if len(code) > 0 {
				return code[0]
			}
			return http.StatusInternalServerError
		}(),
		Data: err,
	})
}
