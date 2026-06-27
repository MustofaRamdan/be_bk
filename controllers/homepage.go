package controllers

import (
	"backend-bk/config"
	"backend-bk/models"

	"github.com/gin-gonic/gin"
)

// alumniStat holds the result of a GROUP BY status query
type alumniStat struct {
	Status string
	Total  int64
}

func GetHomepage(c *gin.Context) {
	// ======================
	// GURU BK
	// ======================
	var guru []models.Guru
	config.DB.
		Order("createdAt desc").
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
	// HERO ARTIKEL
	// ======================
	var hero interface{} = nil
	if len(artikel) > 0 {
		hero = artikel[0]
	}

	// ======================
	// ALUMNI STATS — Single query with GROUP BY
	// Instead of 4 separate COUNT queries, use 1 query
	// ======================
	var stats []alumniStat
	config.DB.Model(&models.Alumni{}).
		Select("status, COUNT(*) as total").
		Where("status_pengajuan = ?", "DITERIMA").
		Group("status").
		Find(&stats)

	var totalAlumni, kuliah, kerja, wirausaha int64
	for _, s := range stats {
		totalAlumni += s.Total
		switch s.Status {
		case "KULIAH":
			kuliah = s.Total
		case "BEKERJA":
			kerja = s.Total
		case "WIRAUSAHA":
			wirausaha = s.Total
		}
	}

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
			"guru":        guru,
			"karya":       karya,
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