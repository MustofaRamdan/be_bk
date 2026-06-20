package controllers

import (
	"backend-bk/config"
	"backend-bk/models"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateGuru(c *gin.Context) {
	var body struct {
		Nama     string `json:"nama"`
		Jabatan  string `json:"jabatan"`
		Kelas    string `json:"kelas"`
		Foto     string `json:"foto"`
		Nip      string `json:"nip"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Body tidak valid"})
		return
	}

	if strings.TrimSpace(body.Nama) == "" {
		c.JSON(400, gin.H{"error": "Nama wajib diisi"})
		return
	}

	if strings.TrimSpace(body.Jabatan) == "" {
		c.JSON(400, gin.H{"error": "Jabatan wajib diisi"})
		return
	}

	if strings.TrimSpace(body.Kelas) == "" {
		c.JSON(400, gin.H{"error": "Kelas wajib dipilih"})
		return
	}

	if strings.TrimSpace(body.Nip) == "" {
		c.JSON(400, gin.H{"error": "NIP wajib diisi"})
		return
	}

	if strings.TrimSpace(body.Password) == "" {
		c.JSON(400, gin.H{"error": "Password wajib diisi"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal memproses password"})
		return
	}

	guru := models.Guru{
		Nama:     body.Nama,
		Jabatan:  body.Jabatan,
		Kelas:    body.Kelas,
		Foto:     &body.Foto,
		Nip:      body.Nip,
		Password: string(hash),
	}

	if err := config.DB.Create(&guru).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal menyimpan data guru"})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"data": guru,
	})
}

func GetGuru(c *gin.Context) {
	var guru []models.Guru

	if err := config.DB.Order("createdAt desc").Find(&guru).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data guru"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": guru,
	})
}

func GetGuruByID(c *gin.Context) {
	id := c.Param("id")

	var guru models.Guru

	if err := config.DB.First(&guru, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Guru tidak ditemukan"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": guru,
	})
}

func UpdateGuru(c *gin.Context) {
	id := c.Param("id")

	var guru models.Guru

	if err := config.DB.First(&guru, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Guru tidak ditemukan"})
		return
	}

	var body map[string]interface{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Body tidak valid"})
		return
	}

	// jika foto diganti, hapus lama
	if newFoto, ok := body["foto"]; ok {
		if guru.Foto != nil && *guru.Foto != "" {
			if newPath, ok := newFoto.(string); ok {
				if newPath != *guru.Foto {
					os.Remove("." + *guru.Foto)
				}
			}
		}
	}

	if passwordVal, ok := body["password"]; ok {
		if passwordStr, ok := passwordVal.(string); ok && strings.TrimSpace(passwordStr) != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(passwordStr), 14)
			if err != nil {
				c.JSON(500, gin.H{"error": "Gagal memproses password"})
				return
			}
			body["password"] = string(hash)
		} else {
			delete(body, "password")
		}
	}

	if err := config.DB.Model(&guru).Updates(body).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal update guru"})
		return
	}

	config.DB.First(&guru, id)

	c.JSON(200, gin.H{
		"success": true,
		"data": guru,
	})
}

func DeleteGuru(c *gin.Context) {
	id := c.Param("id")

	var guru models.Guru

	if err := config.DB.First(&guru, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Guru tidak ditemukan"})
		return
	}

	// hapus foto
	if guru.Foto != nil && *guru.Foto != "" {
		os.Remove("." + *guru.Foto)
	}

	config.DB.Delete(&guru)

	c.JSON(200, gin.H{
		"success": true,
	})
}