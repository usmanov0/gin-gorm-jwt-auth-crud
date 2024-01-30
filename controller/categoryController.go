package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"net/http"
	"simple-crud-api/models"
	"simple-crud-api/pkg/errors"
	"simple-crud-api/pkg/pagination"
	"simple-crud-api/pkg/util"
	"simple-crud-api/storage/initializers"
	"strconv"
)

func CreateCategory(c *gin.Context) {
	var category struct {
		Name string `json:"name" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"validation": util.FormatValidationErrors(errs),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if util.IsUniqueValue("categories", "name", category.Name) ||
		util.IsUniqueValue("categories", "slug", slug.Make(category.Name)) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				"Name": "The name is already exist!",
			},
		})
		return
	}

	categoryModel := models.Category{
		Name: category.Name,
	}

	result := initializers.DB.Create(&categoryModel)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't create category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": categoryModel,
	})
}

func GetCategories(c *gin.Context) {
	var categories []models.Category

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("page", "5")
	perPage, _ := strconv.Atoi(perPageStr)

	res, err := pagination.Paginate(initializers.DB, page, perPage, nil, &categories)
	if err != nil {
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res,
	})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var category struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validation": util.FormatValidationErrors(errs),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var categoryModel models.Category

	result := initializers.DB.First(&categoryModel, id)

	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	if (categoryModel.Name != category.Name &&
		util.IsUniqueValue("categories", "name", category.Name)) ||
		(categoryModel.Name != category.Name &&
			util.IsUniqueValue("categories", "slug", slug.Make(category.Name))) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				" Name": "The name is already exist!",
			},
		})

		return
	}

	updateCategory := models.Category{
		Name: category.Name,
		Slug: slug.Make(category.Name),
	}

	result = initializers.DB.Model(&categoryModel).Updates(updateCategory)
	if err := result.Error; err != nil {
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": updateCategory,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	result := initializers.DB.First(&category, id)
	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)

	}

	initializers.DB.Delete(&category)

	c.JSON(http.StatusOK, gin.H{
		"message": "category deleted successfully",
	})
}
