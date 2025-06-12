package main

import (
	"github.com/gin-gonic/gin"
	auth "github.com/twoonefour/alist-auth"
	"github.com/twoonefour/alist-auth/utils"
)

func main() {
	r := gin.New()
	r.Use(utils.LoggerMiddleware())
	api := r.Group("/alist")
	auth.Setup(api)
	r.Run(":3002")
}
