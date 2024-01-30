package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"simple-crud-api/models"
	"simple-crud-api/pkg/errors"
	"simple-crud-api/pkg/helper"
	"simple-crud-api/pkg/util"
	"simple-crud-api/storage/initializers"
)

func CommentOnPost(c *gin.Context) {
	var comment struct {
		PostId uint   `json:"postId" binding:"required,min=1"`
		Body   string `json:"body" binding:"required,min=1"`
	}

	err := c.ShouldBindJSON(&comment)

	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": util.FormatValidationErrors(errors),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !util.IsExistValue("posts", "id", comment.PostId) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"PostId": "The post does not exist",
			},
		})
		return
	}

	authUser, _ := helper.GetAuthUser(c)

	commentModel := models.Comment{
		PostId: comment.PostId,
		Body:   comment.Body,
		UserId: authUser.Id,
	}

	res := initializers.DB.Create(&commentModel)

	if res.Error != nil {
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": commentModel,
	})
}

func UpdateComment(c *gin.Context) {
	id := c.Param("id")

	var comment struct {
		Body string `json:"body" binding:"required,min=1"`
	}

	err := c.ShouldBindJSON(&comment)

	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": util.FormatValidationErrors(errors),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var commentModel models.Comment
	result := initializers.DB.First(&commentModel, id)

	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	commentModel.Body = comment.Body
	result = initializers.DB.Save(&comment)

	if result.Error != nil {
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

func DeleteComment(c *gin.Context) {
	id := c.Param("comment_id")

	var comment models.Comment
	result := initializers.DB.First(&comment, id)

	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
		return
	}

	initializers.DB.Delete(&comment)

	c.JSON(http.StatusOK, gin.H{
		"message": "The comment has been deleted successfully!",
	})
}
