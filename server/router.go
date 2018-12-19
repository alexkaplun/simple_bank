package server

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"simple_bank/handlers"
)

func NewRouter() *gin.Engine {
	gin.DisableConsoleColor()

	f, _ := os.Create("log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())


	router.GET("/", func(c *gin.Context) {
		c.String(200, "This is your banking application")
	})

	router.PUT("/createAccount", handlers.CreateAccountHandler)
	router.GET("/balance/:id", handlers.GetBalanceByIdHandler)
	router.POST("/transfer", handlers.TransferHandler)
	return router
}
