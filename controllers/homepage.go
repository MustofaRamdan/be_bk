package controllers

import (
	"backend-bk/config"
	"backend-bk/models"

	"github.com/gin-gonic/gin"
)

func GetHomepage(c *gin.Context) {
	// ======================
	// GURU BK
	// ======================
	var guru []models.Guru
	config.DB.
		Order("createdAt desc"). // atau "createdAt" sesuai nama kolom DB Anda
		Limit(5).
		Find(&guru)

	// ======================
	// KARYA SISWA
	// ======================
	var karya []models.Karya
	config.DB.
		Where("status = ?", "DITERIMA").
		Order("createdAt desc").
		Limit(4).
		Find(&karya)

	// ======================
	// ARTIKEL TERBARU
	// ======================
	var artikel []models.Post
	err := config.DB.
		Where("published = ?", true).
		Order("createdAt desc").
		Limit(4).
		Find(&artikel).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal mengambil artikel",
		})
		return
	}

	// ======================
	// HERO ARTI	EL
	// ======================
	var hero interface{} = nil
	if len(artikel) > 0 {
		hero = artikel[0]
	}

	// ======================
	// TOTAL ALUMNI
	// ======================
	var totalAlumni int64
	config.DB.
		Model(&models.Alumni{}).
		Where("status_pengajuan = ?", "DITERIMA").
		Count(&totalAlumni)

	// ======================
	// STATS KARIR ALUMNI
	// ======================
	var kuliah, kerja, wirausaha int64
	config.DB.Model(&models.Alumni{}).
		Where("status_pengajuan = ? AND status = ?", "DITERIMA", "KULIAH").
		Count(&kuliah)
	config.DB.Model(&models.Alumni{}).
		Where("status_pengajuan = ? AND status = ?", "DITERIMA", "BEKERJA").
		Count(&kerja)
	config.DB.Model(&models.Alumni{}).
		Where("status_pengajuan = ? AND status = ?", "DITERIMA", "WIRAUSAHA").
		Count(&wirausaha)

	var persenKuliah, persenKerja, persenWirausaha int64
	if totalAlumni > 0 {
		persenKuliah = (kuliah * 100) / totalAlumni
		persenKerja = (kerja * 100) / totalAlumni
		persenWirausaha = (wirausaha * 100) / totalAlumni
	}

	// ======================
	// RESPONSE
	// ======================
	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"hero":        hero,
			"artikel":     artikel,
			"guru":        guru,    // ← TAMBAH
			"karya":       karya,   // ← TAMBAH
			"totalAlumni": totalAlumni,
			"statsKarir": gin.H{
				"kuliah": gin.H{
					"persen": persenKuliah,
					"jumlah": kuliah,
				},
				"kerja": gin.H{
					"persen": persenKerja,
					"jumlah": kerja,
				},
				"wirausaha": gin.H{
					"persen": persenWirausaha,
					"jumlah": wirausaha,
				},
			},
		},
	})
}