package controllers
import "fmt"

import (
	"backend-bk/config"
	"backend-bk/models"
	"backend-bk/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	fmt.Println("REGISTER HIT")
	var input models.User
	fmt.Println(input)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"message": "bind gagal",
			"error": err.Error(),
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	user := models.User{
		Name: input.Name,
		Email: input.Email,
		Password: string(hash),
	}

	result := config.DB.Debug().Create(&user)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "insert gagal",
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "berhasil",
		"user": user,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Nip      string `json:"nip"`
		Password string `json:"password"`
	}
	var guru models.Guru

	c.ShouldBindJSON(&input)

	config.DB.Where("nip = ?", input.Nip).First(&guru)

	if guru.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "NIP salah"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(guru.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Password salah"})
		return
	}

	token, _ := utils.GenerateToken(guru.ID, guru.Nip)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token": token,
		"user": gin.H{
			"id":   guru.ID,
			"nama": guru.Nama,
			"nip":  guru.Nip,
		},
	})
}