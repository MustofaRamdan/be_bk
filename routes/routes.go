package routes

import (
	"backend-bk/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	api.GET("/posts", controllers.GetPosts)
	api.POST("/posts", controllers.CreatePost)
	api.POST("/upload", controllers.Upload)
	api.GET("/posts/:id", controllers.GetPostByID)
	api.PUT("/posts/:id", controllers.UpdatePost)
	api.DELETE("/posts/:id", controllers.DeletePost)

	api.POST("/karya", controllers.CreateKarya)
	api.GET("/karya", controllers.GetKarya)
	api.GET("/karya/:id", controllers.GetKaryaByID)
	api.PUT("/karya/:id", controllers.UpdateKarya)
	api.DELETE("/karya/:id", controllers.DeleteKarya)

	api.POST("/guru", controllers.CreateGuru)
	api.GET("/guru", controllers.GetGuru)
	api.GET("/guru/:id", controllers.GetGuruByID)
	api.PUT("/guru/:id", controllers.UpdateGuru)
	api.DELETE("/guru/:id", controllers.DeleteGuru)

	api.POST("/alumni", controllers.CreateAlumni)
	api.GET("/alumni", controllers.GetAlumni)
	api.GET("/alumni/:id", controllers.GetAlumniByID)
	api.PUT("/alumni/:id", controllers.UpdateAlumniStatus)

	api.GET("/konseling", controllers.GetKonseling)
	api.GET("/konseling/:id", controllers.GetKonselingByID)
	api.PUT("/konseling/:id", controllers.UpdateJawabanKonseling)
	
	api.POST("/konseling", controllers.CreateKonseling)
	api.GET("/dashboard", controllers.GetDashboard)

	api.GET("/homepage", controllers.GetHomepage)
}