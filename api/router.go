package api

import (
	"github.com/gin-gonic/gin"
	//swaggerfile "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	"simple-crud-api/controller"
	"simple-crud-api/middleware"
)

func Route() *gin.Engine {
	r := gin.Default()
	r.POST("/api/sign-up", controller.SignUp)
	r.POST("/api/log-in", controller.SignIn)

	r.Use(middleware.RequireAuth)
	r.POST("/api/log-out", controller.LogOut)
	userRouter := r.Group("/api/users")
	{
		userRouter.GET("/", controller.GetUsers)
		userRouter.PUT("/update", controller.UpdateUser)
		userRouter.DELETE("/delete", controller.DeleteUser)
	}

	categoryRouter := r.Group("/api/categories")
	{
		categoryRouter.POST("/create", controller.CreateCategory)
		categoryRouter.GET("/", controller.GetCategories)
		categoryRouter.PUT("/update", controller.UpdateCategory)
		categoryRouter.DELETE("/delete", controller.DeleteCategory)
	}

	postRouter := r.Group("/api/posts")
	{
		postRouter.POST("/create", controller.CreatePost)
		postRouter.GET("/", controller.GetPosts)
		postRouter.GET("/read-post", controller.ReadPosts)
		postRouter.GET("/edit", controller.EditPost)
		postRouter.PUT("/update", controller.UpdatePost)
		postRouter.DELETE("/delete", controller.DeletePost)
	}

	commentRouter := r.Group("/api/comments")
	{
		commentRouter.POST("/comment", controller.CommentOnPost)
		commentRouter.PUT("/update", controller.UpdateComment)
		commentRouter.DELETE("/delete", controller.DeleteComment)
	}
	return r
}
