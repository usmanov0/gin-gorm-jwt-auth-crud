package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"net/http"
	"simple-crud-api/pkg/errors"
	"simple-crud-api/pkg/helper"
	"simple-crud-api/pkg/pagination"
	"simple-crud-api/pkg/util"
	"simple-crud-api/storage/initializers"
	"strconv"
)

type Category struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Posts []Post `json:"posts"`
}

// @Summary Create a new category
// @Description Create a new category
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param name body string true "Category name"
// @Success 200 {object} Category
// @Failure 400 {object}
// @Failure 401 {object}
// @Failure 500 {object}
// @Router /api/categories [post]
func CreateCategory(c *gin.Context) {
	_, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	categoryModel := Category{
		Name: category.Name,
		Slug: slug.Make(category.Name),
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

// @Summary Get a list of categories
// @Description Get a list of categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object}  pagination.PaginateRes
// @Failure 401
// @Failure 500
// @Router /api/categories/ [get]
func GetCategories(c *gin.Context) {
	var categories []Category

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

// @Summary Update a category
// @Description Update a category
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param id path int true "Category ID"
// @Param name body string true "Category name"
// @Success 200 {object} Category "Updated category"
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /api/categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	_, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	var categoryModel Category

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

	updateCategory := Category{
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

// @Summary Delete a category
// @Description Delete a category
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param id path int true "Category ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /api/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	_, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	var category Category

	result := initializers.DB.First(&category, id)
	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)

	}

	initializers.DB.Delete(&category)

	c.JSON(http.StatusOK, gin.H{
		"message": "category deleted successfully",
	})
}
