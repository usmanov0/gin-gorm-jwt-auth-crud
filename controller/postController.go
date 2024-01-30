package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"os/user"
	"simple-crud-api/pkg/errors"
	"simple-crud-api/pkg/helper"
	"simple-crud-api/pkg/pagination"
	"simple-crud-api/pkg/util"
	"simple-crud-api/storage/initializers"
	"strconv"
)

type Post struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	UserId     uint      `json:"user_id"`
	CategoryId uint      `json:"category_id"`
	Category   Category  `json:"category"`
	User       user.User `json:"user"`
	Comments   []Comment `json:"comments"`
}

// @Summary Create a new post
// @Description Create a new post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
//
//	@Param post body struct {
//	  Title      string `json:"title" binding:"required,min=2,max=200"`
//	  Body       string `json:"body" binding:"required"`
//	  CategoryId uint   `json:"categoryId" binding:"required,min=1"`
//	} true "Post details"
//
// @Success 200 {object} "Newly created post"
// @Failure 401 {object} "Unauthorized"
// @Router /api/posts [post]
func CreatePost(c *gin.Context) {
	authUser, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	postModel := Post{
		Title:      post.Title,
		Body:       post.Body,
		CategoryId: post.CategoryId,
		UserId:     authUser.Id,
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

// @Summary Get a list of posts
// @Description Get a list of posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param page query int false "Page number"
// @Param perPage query int false "Number of items per page"
// @Success 200 {object} pagination.PaginateRes
// @Failure 401 {object} "Unauthorized"
// @Failure 500 {object} "Internal Server Error"
// @Router /api/posts [get]
func GetPosts(c *gin.Context) {
	var posts []Post

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

// @Summary Read a post by ID
// @Description Read a post by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param id path int true "Post ID"
// @Success 200 {object} Post
// @Failure 401 {object} "Unauthorized"
// @Failure 404 {object} "Post not found"
// @Router /api/posts/read-post [get]
func ReadPosts(c *gin.Context) {
	id := c.Param("id")

	var post Post
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

// @Summary Read a post by ID
// @Description Read a post by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param id path int true "Post ID"
// @Success 200 {object} "Post details"
// @Failure 401 {object} "Unauthorized"
// @Failure 404 {object} "Post not found"
// @Router /api/posts/edit/{id} [get]
func EditPost(c *gin.Context) {
	authUser, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")

	var post Post
	if post.UserId != authUser.Id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Forbidden: You are not allowed to edit this post"})
		return
	}

	result := initializers.DB.First(&post, id)

	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// @Summary Update a post by ID
// @Description Update a post by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param id path int true "Post ID"
//
//	@Param post body struct {
//	  Title      string `json:"title" binding:"required,min=2,max=200"`
//	  Body       string `json:"body" binding:"required"`
//	  CategoryId uint   `json:"categoryId" binding:"required,min=1"`
//	} true "Updated post details"
//
// @Success 200 {object} Post
// @Failure 401 {object} "Unauthorized"
// @Failure 403 {object} "Forbidden"
// @Router /api/posts/update/{id} [put]
func UpdatePost(c *gin.Context) {
	authUser, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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

	var postModel Post
	result := initializers.DB.First(&postModel, id)
	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	if postModel.UserId != authUser.Id {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: You are not allowed to update this post",
		})
		return
	}

	updatePost := Post{
		Title:      post.Title,
		Body:       post.Body,
		CategoryId: post.CategoryId,
		UserId:     authUser.Id,
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

// @Summary Delete a post by ID
// @Description Delete a post by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param id path int true "Post ID"
// @Success 200 {object} "Post deleted successfully"
// @Failure 401 {object} "Unauthorized"
// @Failure 403 {object} "Forbidden"
// @Failure 404 {object} "Post not found"
// @Router /api/posts/delete/{id} [delete]
func DeletePost(c *gin.Context) {
	authUser, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	var post Post

	result := initializers.DB.First(&post, id)
	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	if post.UserId != authUser.Id {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: You are not allowed to delete this post",
		})
		return
	}

	initializers.DB.Delete(&post)

	c.JSON(http.StatusOK, gin.H{
		"message": "post deleted successfully",
	})
}
