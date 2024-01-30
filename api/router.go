package api

import (
	"github.com/gin-gonic/gin"
	//swaggerfile "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	"simple-crud-api/controller"
	"simple-crud-api/middleware"
)

func Route(r *gin.Engine) {
	r.POST("/api/sign-up", controller.SignUp)
	r.POST("/api/log-in", controller.SignIn)

	r.Use(middleware.RequireAuth)
	r.POST("/api/log-out", controller.LogOut)
	userRouter := r.Group("/api/users")
	{
		userRouter.GET("/", controller.GetUsers)
		userRouter.PUT("/update/:id", controller.UpdateUser)
		userRouter.DELETE("/delete/:id", controller.DeleteUser)
	}

	categoryRouter := r.Group("/api/categories")
	{
		categoryRouter.POST("/create", controller.CreateCategory)
		categoryRouter.GET("/", controller.GetCategories)
		categoryRouter.PUT("/update/:id", controller.UpdateCategory)
		categoryRouter.DELETE("/delete/:id", controller.DeleteCategory)
	}

	postRouter := r.Group("/api/posts")
	{
		postRouter.POST("/create", controller.CreatePost)
		postRouter.GET("/", controller.GetPosts)
		postRouter.GET("/read-post/:id", controller.ReadPosts)
		postRouter.GET("/edit/:id", controller.EditPost)
		postRouter.PUT("/update/:id", controller.UpdatePost)
		postRouter.DELETE("/delete/:id", controller.DeletePost)
	}

	commentRouter := r.Group("/api/comments")
	{
		commentRouter.POST("/comment", controller.CommentOnPost)
		commentRouter.PUT("/update/:id", controller.UpdateComment)
		commentRouter.DELETE("/delete/:id", controller.DeleteComment)
	}
}
