package main



import (
	"backend-bk/config"
	"backend-bk/routes"
	"backend-bk/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"time"
    "github.com/gin-contrib/cors"
)

func main() {
	godotenv.Load()

	config.ConnectDB()
    config.DB.AutoMigrate(
	&models.Post{},
	&models.User{},
	&models.Guru{},
	&models.Karya{},
	&models.Alumni{},
	&models.Konseling{},
)

	r := gin.Default()
    r.Static("/uploads", "./uploads")

	r.Use(cors.New(cors.Config{
    AllowOrigins: []string{
        "http://localhost:3000",
    },
    AllowMethods: []string{
        "GET", "POST", "PUT", "DELETE", "OPTIONS",
    },
    AllowHeaders: []string{
        "Origin", "Content-Type", "Authorization",
    },
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
}))

	routes.SetupRoutes(r)
    

	r.Run(":3001")
}