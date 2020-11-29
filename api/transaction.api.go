package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//SetupTransactionApi - send to api.go
func SetupTransactionApi(r *gin.Engine) {
	transactionApi := r.Group("/api/v2")
	{
		transactionApi.GET("/transaction", getTransaction)
		transactionApi.POST("/transaction", createTransaction)
	}
}

func getTransaction(c *gin.Context) {
	c.JSON(http.StatusOK, "List Transaction")
}

func createTransaction(c *gin.Context) {
	c.JSON(http.StatusOK, "List Transaction")
}
