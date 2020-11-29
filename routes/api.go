package routes

import (
	"app/controllers"
	"app/db"
	"app/interceptor"

	"github.com/gin-gonic/gin"
)

//Setup - send function Setup
func Setup(r *gin.Engine) {
	db.SetupDB()

	authenAPI := r.Group("api/v2")
	{
		authenAPI.POST("/login", controllers.Login)
		authenAPI.POST("/register", controllers.Register)
	}

	productAPI := r.Group("/api/v2")
	{
		productAPI.GET("/products", controllers.GetAllProduct)
		productAPI.GET("/product", controllers.GetProduct)
		productAPI.GET("/product/:id", controllers.GetProductByID)
		productAPI.POST("/product" /*interceptor.JwtVerify,*/, controllers.CreateProduct)
		productAPI.PUT("/product" /*interceptor.JwtVerify,*/, controllers.EditProduct)
		productAPI.DELETE("/product/:id" /*interceptor.JwtVerify,*/, controllers.DeleteProduct)
	}

	transactionApi := r.Group("/api/v2")
	{
		transactionApi.GET("/transaction", controllers.GetTransaction)
		transactionApi.POST("/transaction", interceptor.JwtVerify, controllers.CreateTransaction)
		transactionApi.PATCH("/transaction/:id", controllers.UpdateTransaction)
		transactionApi.DELETE("/transaction/:id", controllers.DeleteTransaction)
	}
}
