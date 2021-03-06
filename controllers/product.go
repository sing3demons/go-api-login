package controllers

import (
	"app/db"
	"strconv"

	"app/models"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type productResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Stock int64   `json:"stock"`
	Price float64 `json:"price"`
	Image string  `json:"image"`
}

type createProductForm struct {
	Name  string                `form:"name" binding:"required"`
	Stock int64                 `form:"stock" binding:"required"`
	Price float64               `form:"price" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

type updateProductForm struct {
	ID    uint                  `form:"id"`
	Name  string                `form:"name"`
	Stock int64                 `form:"stock"`
	Price float64               `form:"price"`
	Image *multipart.FileHeader `form:"image"`
}

func GetProduct(ctx *gin.Context) {
	var product models.Product
	keyword := ctx.Query("keyword")
	if keyword != "" {
		keyword = fmt.Sprintf("%%%s%%", keyword)
		db.GetDB().Where("name like ?", keyword).Find(&product)
		return
	}
	db.GetDB().Find(&product)
	ctx.JSON(200, product)
}

func GetProductByID(c *gin.Context) {
	var product *models.Product
	var err error
	product, err = findByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var serializedProduct productResponse
	copier.Copy(&serializedProduct, &product)
	c.JSON(http.StatusOK, gin.H{"product": serializedProduct})
}

func findByID(ctx *gin.Context) (*models.Product, error) {
	var product models.Product
	id := ctx.Param("id")

	if err := db.GetDB().First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func GetAllProduct(c *gin.Context) {
	var products []models.Product
	if err := db.GetDB().Order("id desc").Find(&products).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	var serializedProducts []productResponse
	copier.Copy(&serializedProducts, &products)
	c.JSON(200, gin.H{"result": serializedProducts})
}

func CreateProduct(c *gin.Context) {
	var form createProductForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	copier.Copy(&product, &form)

	if err := db.GetDB().Create(&product).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	image, _ := c.FormFile("image")
	saveImage(image, &product, c)

	c.JSON(http.StatusOK, gin.H{"result": product})

}

//EditProduct
func EditProduct(c *gin.Context) {
	var product models.Product

	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 32)
	product.ID = uint(id)
	product.Name = c.PostForm("name")
	product.Stock, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64)
	product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)

	if err := db.GetDB().Save(&product).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	saveImage(image, &product, c)

	serializedProduct := productResponse{}
	copier.Copy(&serializedProduct, &product)
	c.JSON(http.StatusOK, gin.H{"product": serializedProduct})
}

//Delete/:id
func DeleteProduct(ctx *gin.Context) {
	product, err := findByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := db.GetDB().Unscoped().Delete(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func saveImage(image *multipart.FileHeader, product *models.Product, c *gin.Context) {
	if image != nil {
		runningDir, _ := os.Getwd()
		product.Image = image.Filename
		extension := filepath.Ext(image.Filename)
		fileName := fmt.Sprintf("%d%s", product.ID, extension)
		filePath := fmt.Sprintf("%s/uploaded/images/%s", runningDir, fileName)

		if fileExists(filePath) {
			os.Remove(filePath)
		}
		c.SaveUploadedFile(image, filePath)
		db.GetDB().Model(&product).Update("image", fileName)
	}
}
