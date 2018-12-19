package server

import (
	"github.com/gin-gonic/gin"
	"simple_bank/handlers"
	"simple_bank/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	//router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "This is your banking application")
	})

	router.PUT("/createAccount", handlers.CreateAccountHandler)
	router.GET("/balance/:id", handlers.GetBalanceByIdHandler)
	router.POST("/transfer", handlers.TransferHandler)

	router.Use(middlewares.AuthMiddleware())
	return router
}
