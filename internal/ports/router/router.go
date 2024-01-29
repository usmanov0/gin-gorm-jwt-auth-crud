package router

import (
	"github.com/gin-gonic/gin"
	"simple-crud-api/internal/middleware"
	handler2 "simple-crud-api/internal/ports/handler"
)

func Route(r *gin.Engine) {
	r.POST("/api/sign-up", handler2.SignUp)
	r.POST("/api/log-in", handler2.SignIn)

	r.Use(middleware.RequireAuth)
	r.POST("/api/log-out", handler2.LogOut)
	userRouter := r.Group("/api/users")
	{
		userRouter.GET("/", handler2.GetUsers)
		userRouter.PUT("/update", handler2.UpdateUser)
		userRouter.DELETE("/delete", handler2.DeleteUser)
	}

	categoryRouter := r.Group("/api/categories")
	{
		categoryRouter.POST("/create", handler2.CreateCategory)
		categoryRouter.GET("/", handler2.GetCategories)
		categoryRouter.PUT("/update", handler2.UpdateCategory)
		categoryRouter.DELETE("/delete", handler2.DeleteCategory)
	}

	postRouter := r.Group("/api/posts")
	{
		postRouter.POST("/create", handler2.CreatePost)
		postRouter.GET("/", handler2.GetPosts)
		postRouter.GET("/read-post", handler2.ReadPosts)
		postRouter.GET("/edit", handler2.EditPost)
		postRouter.PUT("/update", handler2.UpdatePost)
		postRouter.DELETE("/delete", handler2.DeletePost)
	}

	commentRouter := r.Group("/api/comment")
	{
		commentRouter.POST("/comment", handler2.CommentOnPost)
		commentRouter.PUT("/update", handler2.UpdateComment)
		commentRouter.DELETE("/delete", handler2.DeleteComment)
	}

}
