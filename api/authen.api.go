package api

import (
	"app/db"
	"app/interceptor"
	"app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type formLogin struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type userResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Level    string
}

//SetupAuthenAPI - set-apisendTOAPISetupAuthenAPI
func SetupAuthenAPI(r *gin.Engine) {
	authenAPI := r.Group("api/v2")
	{
		authenAPI.POST("/login", login)
		authenAPI.POST("/register", register)
	}

}

func login(c *gin.Context) {
	var form formLogin
	var user models.User
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(401, gin.H{"status": "unable to bind data"})
		return
	}

	if err := db.GetDB().Where("username = ?", form.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "nok", "error": err.Error()})
		return
	}

	if checkPasswordHash(form.Password, user.Password) == false {
		c.JSON(401, gin.H{"result": "nok", "error": "invalid password"})
		return
	}

	serializedUser := interceptor.JwtSign(user)
	// copier.Copy(&form, &user)
	c.JSON(http.StatusOK, gin.H{"data": serializedUser})
}

func register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(401, gin.H{"status": "unable to bind data"})
		return
	}

	user.Password, _ = hashPassword(user.Password)
	if err := db.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"result": "nok", "error": err})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	c.JSON(http.StatusOK, gin.H{"result": "ok", "data": serializedUser})

}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
