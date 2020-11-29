package controllers

import (
	"app/db"
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

type transactionResult struct {
	ID            uint
	Total         float64
	Paid          float64
	Change        float64
	PaymentType   string
	PaymentDetail string
	OrderList     string
	Staff         string
}

func GetTransaction(c *gin.Context) {
	var result []transactionResult
	db.GetDB().Debug().Raw("SELECT transactions.id, total, paid, change, payment_type, payment_detail, order_list, users.username as Staff, transactions.created_at FROM transactions join users on transactions.staff_id = users.id", nil).Scan(&result)

	c.JSON(http.StatusOK, gin.H{"trsnsactions": result})
}

func CreateTransaction(c *gin.Context) {
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

func UpdateTransaction(c *gin.Context) {
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

func DeleteTransaction(c *gin.Context) {
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
