package routes

import (
	"backend-bk/controllers"
	"backend-bk/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Public Routes
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
	api.GET("/homepage", controllers.GetHomepage)
	
	api.GET("/posts", controllers.GetPosts)
	api.GET("/posts/:id", controllers.GetPostByID)
	
	api.GET("/karya", controllers.GetKarya)
	api.GET("/karya/:id", controllers.GetKaryaByID)
	
	api.GET("/guru", controllers.GetGuru)
	api.GET("/guru/:id", controllers.GetGuruByID)
	
	api.POST("/alumni", controllers.CreateAlumni)
	api.POST("/konseling", controllers.CreateKonseling)
	api.POST("/upload", controllers.Upload)

	// Protected Routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/posts", controllers.CreatePost)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)

		protected.POST("/karya", controllers.CreateKarya)
		protected.PUT("/karya/:id", controllers.UpdateKarya)
		protected.DELETE("/karya/:id", controllers.DeleteKarya)

		protected.POST("/guru", controllers.CreateGuru)
		protected.PUT("/guru/:id", controllers.UpdateGuru)
		protected.DELETE("/guru/:id", controllers.DeleteGuru)

		protected.GET("/alumni", controllers.GetAlumni)
		protected.GET("/alumni/:id", controllers.GetAlumniByID)
		protected.PUT("/alumni/:id", controllers.UpdateAlumniStatus)

		protected.GET("/konseling", controllers.GetKonseling)
		protected.GET("/konseling/:id", controllers.GetKonselingByID)
		protected.PUT("/konseling/:id", controllers.UpdateJawabanKonseling)

		protected.GET("/dashboard", controllers.GetDashboard)
	}
}