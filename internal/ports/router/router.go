package router

import (
	"github.com/gin-gonic/gin"
	"simple-crud-api/internal/middleware"
	handler "simple-crud-api/internal/ports/handler"
)

func Route(r *gin.Engine) {
	r.POST("/api/sign-up", handler.SignUp)
	r.POST("/api/log-in", handler.SignIn)

	r.Use(middleware.RequireAuth)
	r.POST("/api/log-out", handler.LogOut)
	userRouter := r.Group("/api/users")
	{
		userRouter.GET("/", handler.GetUsers)
		userRouter.PUT("/update", handler.UpdateUser)
		userRouter.DELETE("/delete", handler.DeleteUser)
	}

	categoryRouter := r.Group("/api/categories")
	{
		categoryRouter.POST("/create", handler.CreateCategory)
		categoryRouter.GET("/", handler.GetCategories)
		categoryRouter.PUT("/update", handler.UpdateCategory)
		categoryRouter.DELETE("/delete", handler.DeleteCategory)
	}

	postRouter := r.Group("/api/posts")
	{
		postRouter.POST("/create", handler.CreatePost)
		postRouter.GET("/", handler.GetPosts)
		postRouter.GET("/read-post", handler.ReadPosts)
		postRouter.GET("/edit", handler.EditPost)
		postRouter.PUT("/update", handler.UpdatePost)
		postRouter.DELETE("/delete", handler.DeletePost)
	}

	commentRouter := r.Group("/api/comments")
	{
		commentRouter.POST("/comment", handler.CommentOnPost)
		commentRouter.PUT("/update", handler.UpdateComment)
		commentRouter.DELETE("/delete", handler.DeleteComment)
	}

}
