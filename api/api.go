package api

import (
	"app/db"

	"github.com/gin-gonic/gin"
)

//Setup - send function Setup
func Setup(r *gin.Engine) {
	db.SetupDB()
	SetupAuthenAPI(r)
	SetupProductAPI(r)
	SetupTransactionApi(r)
}
