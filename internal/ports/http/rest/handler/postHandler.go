package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"simple-crud-api/internal/common/db/initializers"
	"simple-crud-api/internal/errors"
	"simple-crud-api/internal/helper"
	"simple-crud-api/internal/models"
	"simple-crud-api/internal/pagination"
	"simple-crud-api/internal/util"
	"strconv"
)

func CreatePost(c *gin.Context) {
	var post struct {
		Title      string `json:"title" binding:"required,min=2,max=200"`
		Body       string `json:"body" binding:"required"`
		CategoryId uint   `json:"categoryId" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validation": util.FormatValidationErrors(errs),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if !util.IsExistValue("categories", "id", post.CategoryId) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"category id": "The category doesn't exists!",
			},
		})
		return
	}

	authId := helper.GetAuthUser(c).Id

	postModel := models.Post{
		Title:      post.Title,
		Body:       post.Body,
		CategoryId: post.CategoryId,
		UserId:     authId,
	}

	res := initializers.DB.Create(&postModel)

	if res.Error != nil {
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": postModel,
	})
}

func GetPosts(c *gin.Context) {
	var posts []models.Post

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, _ := strconv.Atoi(perPageStr)

	preLoadFunc := func(query *gorm.DB) *gorm.DB {
		return query.Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug").Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			})
		})
	}

	result, err := pagination.Paginate(initializers.DB, page, perPage, preLoadFunc, &posts)

	if err != nil {
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": result,
	})
}

func ReadPosts(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	result := initializers.DB.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, slug")
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Select("id, post_id, user_id, body, created_at")
	}).First(&post, id)

	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func EditPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	result := initializers.DB.First(&post, id)

	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var post struct {
		Title      string `json:"title" binding:"required,min=2,max=200"`
		Body       string `json:"body" binding:"required"`
		CategoryId uint   `json:"categoryId" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": util.FormatValidationErrors(errs),
			})

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var postModel models.Post
	result := initializers.DB.First(&postModel, id)
	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	authId := helper.GetAuthUser(c).Id
	updatePost := models.Post{
		Title:      post.Title,
		Body:       post.Body,
		CategoryId: post.CategoryId,
		UserId:     authId,
	}

	result = initializers.DB.Model(&postModel).Updates(&updatePost)

	if result.Error != nil {
		errors.InternalServerError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"post": updatePost,
	})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	result := initializers.DB.First(&post, id)
	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	initializers.DB.Delete(&post)

	c.JSON(http.StatusOK, gin.H{
		"message": "post deleted successfully",
	})
}
