package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"simple-crud-api/pkg/errors"
	"simple-crud-api/pkg/helper"
	"simple-crud-api/pkg/util"
	"simple-crud-api/storage/initializers"
)

type Comment struct {
	ID     uint   `json:"id"`
	Body   string `json:"body"`
	PostId uint   `json:"post_id" binding:"required, gt=0"`
	UserId uint   `json:"user_id"`
	User   User   `json:"user"`
}

// @Summary Comment on a post
// @Description Comment on a post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
//
//	@Param comment body struct {
//	  PostId uint   `json:"postId" binding:"required,min=1"`
//	  Body   string `json:"body" binding:"required,min=1"`
//	} true "Comment details"
//
// @Success 200 {object} Comment
// @Failure 401 {object}
// @Failure 422 {object}
// @Failure 500 {object}
// @Router /api/comments/comment [post]
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

	commentModel := Comment{
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

// @Summary Update a comment by ID
// @Description Update a comment by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param id path int true "Comment ID"
//
//	@Param comment body struct {
//	  Body string `json:"body" binding:"required,min=1"`
//	} true "Updated comment details"
//
// @Success 200
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 422
// @Failure 500
// @Router /api/comments/update{id} [put]
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

	var commentModel Comment
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

// @Summary Delete a comment by ID
// @Description Delete a comment by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <JWT_TOKEN>"
// @Param comment_id path int true "Comment ID"
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /api/comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	id := c.Param("comment_id")

	var comment Comment
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
