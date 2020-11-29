package api

import (
	"app/db"
	"app/interceptor"
	"app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type transactionForm struct {
	Total         float64 `json:"total"`
	Paid          float64 `json:"paid"`
	Change        float64 `json:"change"`
	PaymentType   string  `json:"payment_type"`
	PaymentDetail string  `json:"payment_detail"`
	OrderList     string  `json:"order_list"`
	StaffID       string  `json:"staff_id"`
}

type updateTransactionForm struct {
	Total         float64 `json:"total"`
	Paid          float64 `json:"paid"`
	Change        float64 `json:"change"`
	PaymentType   string  `json:"payment_type"`
	PaymentDetail string  `json:"payment_detail"`
	OrderList     string  `json:"order_list"`
	StaffID       string  `json:"staff_id"`
}

type transactionResponse struct {
	Total         float64 `json:"total"`
	Paid          float64 `json:"paid"`
	Change        float64 `json:"change"`
	PaymentType   string  `json:"payment_type"`
	PaymentDetail string  `json:"payment_detail"`
	OrderList     string  `json:"order_list"`
	StaffID       string  `json:"staff_id"`
}

//SetupTransactionApi - api api.go
func SetupTransactionApi(r *gin.Engine) {
	transactionApi := r.Group("/api/v2")
	{
		transactionApi.GET("/transaction", getTransaction)
		transactionApi.POST("/transaction", interceptor.JwtVerify, createTransaction)
		transactionApi.PATCH("/transaction/:id", updateTransaction)
		transactionApi.DELETE("/transaction/:id", deleteTransaction)
	}
}

func getTransaction(c *gin.Context) {
	var transaction []models.Transaction
	db.GetDB().Find(&transaction)

	serializedArticle := []transactionResponse{}
	copier.Copy(&serializedArticle, &transaction)
	c.JSON(http.StatusOK, gin.H{"trsnsactions": serializedArticle})
}

func createTransaction(c *gin.Context) {
	var form transactionForm

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var transaction models.Transaction
	copier.Copy(&transaction, &form)
	transaction.StaffID = c.GetString("jwt_staff_id")
	if err := db.GetDB().Create(&transaction).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	serializedArticle := transactionResponse{}
	copier.Copy(&serializedArticle, &transaction)
	c.JSON(http.StatusOK, gin.H{"trsnsaction": serializedArticle})
}

func updateTransaction(c *gin.Context) {
	var form updateTransactionForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	transaction, err := findByTransactionID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	db.GetDB().Model(&transaction).Updates(&form)

	serializedArticle := transactionResponse{}
	copier.Copy(&serializedArticle, &transaction)
	c.JSON(http.StatusOK, gin.H{"trsnsaction": serializedArticle})

}

func deleteTransaction(c *gin.Context) {
	transaction, err := findByTransactionID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.GetDB().Unscoped().Delete(&transaction).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func findByTransactionID(ctx *gin.Context) (*models.Transaction, error) {
	var transaction models.Transaction
	id := ctx.Param("id")

	if err := db.GetDB().First(&transaction, id).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}
