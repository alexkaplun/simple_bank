package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"simple_bank/handlers"
	"simple_bank/middlewares"
)

func NewRouter(logger *zap.Logger) *gin.Engine {
	router := gin.New()

	router.Use(middlewares.ZapLogger(logger))
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "This is your banking application")
	})

	router.PUT("/createAccount", handlers.CreateAccountHandler)
	router.GET("/balance/:id", handlers.GetBalanceByIdHandler)
	router.POST("/transfer", handlers.TransferHandler)

	return router
}
