package controllers

import (
	"backend-bk/config"
	"backend-bk/models"
	"strings"
	"os"
	"github.com/gin-gonic/gin"
)

func CreateGuru(c *gin.Context) {
	var body struct {
		Nama    string `json:"nama"`
		Jabatan string `json:"jabatan"`
		Kelas   string `json:"kelas"`
		Foto    string `json:"foto"`
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

	guru := models.Guru{
		Nama:    body.Nama,
		Jabatan: body.Jabatan,
		Kelas:   body.Kelas,
		Foto:    &body.Foto,
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