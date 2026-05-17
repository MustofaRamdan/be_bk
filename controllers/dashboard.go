package controllers

import (
	"backend-bk/config"
	"backend-bk/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	// =====================
	// VARIABEL
	// =====================
	var totalSiswa int64
	var totalArtikel int64
	var totalGuru int64
	var totalKarya int64
	var totalAlumni int64

	var kuliah int64
	var bekerja int64
	var wirausaha int64

	var pendingKarya int64
	var pendingAlumni int64

	now := time.Now()

	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	firstDayOfLastMonth := firstDayOfMonth.AddDate(0, -1, 0)

	// =====================
	// TOTAL DATA
	// =====================

	config.DB.Model(&models.User{}).Count(&totalSiswa)

	config.DB.Model(&models.Post{}).
		Where("published = ?", true).
		Count(&totalArtikel)

	config.DB.Model(&models.Guru{}).Count(&totalGuru)

	config.DB.Model(&models.Karya{}).Count(&totalKarya)

	config.DB.Model(&models.Alumni{}).
		Where("status_pengajuan = ?", "DITERIMA").
		Count(&totalAlumni)

	// =====================
	// STATUS ALUMNI
	// =====================

	config.DB.Model(&models.Alumni{}).
		Where("status = ?", "KULIAH").
		Count(&kuliah)

	config.DB.Model(&models.Alumni{}).
		Where("status = ?", "BEKERJA").
		Count(&bekerja)

	config.DB.Model(&models.Alumni{}).
		Where("status = ?", "WIRAUSAHA").
		Count(&wirausaha)

	// =====================
	// PENDING
	// =====================

	config.DB.Model(&models.Karya{}).
		Where("status = ?", "PENDING").
		Count(&pendingKarya)

	config.DB.Model(&models.Alumni{}).
		Where("status_pengajuan = ?", "PENDING").
		Count(&pendingAlumni)

	// =====================
	// GROWTH ARTIKEL
	// =====================

	var artikelBulanIni int64
	var artikelBulanLalu int64

	config.DB.Model(&models.Post{}).
		Where("published = ? AND createdAt >= ?", true, firstDayOfMonth).
		Count(&artikelBulanIni)

	config.DB.Model(&models.Post{}).
		Where("published = ? AND createdAt >= ? AND createdAt < ?", true, firstDayOfLastMonth, firstDayOfMonth).
		Count(&artikelBulanLalu)

	artikelGrowth := calculateGrowth(artikelBulanIni, artikelBulanLalu)

	// =====================
	// GROWTH GURU
	// =====================

	var guruBulanIni int64
	var guruBulanLalu int64

	config.DB.Model(&models.Guru{}).
		Where("createdAt >= ?", firstDayOfMonth).
		Count(&guruBulanIni)

	config.DB.Model(&models.Guru{}).
		Where("createdAt >= ? AND createdAt < ?", firstDayOfLastMonth, firstDayOfMonth).
		Count(&guruBulanLalu)

	guruGrowth := calculateGrowth(guruBulanIni, guruBulanLalu)

	// =====================
	// GROWTH KARYA
	// =====================

	var karyaBulanIni int64
	var karyaBulanLalu int64

	config.DB.Model(&models.Karya{}).
		Where("createdAt >= ?", firstDayOfMonth).
		Count(&karyaBulanIni)

	config.DB.Model(&models.Karya{}).
		Where("createdAt >= ? AND createdAt < ?", firstDayOfLastMonth, firstDayOfMonth).
		Count(&karyaBulanLalu)

	karyaGrowth := calculateGrowth(karyaBulanIni, karyaBulanLalu)

	// =====================
	// GROWTH ALUMNI (DITERIMA)
	// =====================

	var alumniBulanIni int64
	var alumniBulanLalu int64

	config.DB.Model(&models.Alumni{}).
		Where("status_pengajuan = ? AND created_at >= ?", "DITERIMA", firstDayOfMonth).
		Count(&alumniBulanIni)

	config.DB.Model(&models.Alumni{}).
		Where("status_pengajuan = ? AND created_at >= ? AND created_at < ?", "DITERIMA", firstDayOfLastMonth, firstDayOfMonth).
		Count(&alumniBulanLalu)

	alumniGrowth := calculateGrowth(alumniBulanIni, alumniBulanLalu)

	// =====================
	// AKTIVITAS TERBARU
	// =====================

	var aktivitas []gin.H

	var posts []models.Post
	config.DB.Order("createdAt desc").Limit(3).Find(&posts)

	for _, p := range posts {
		title := p.Title
		if len(title) > 30 {
			title = title[:30] + "..."
		}

		aktivitas = append(aktivitas, gin.H{
			"id":      p.ID,
			"nama":    "Admin",
			"aksi":    "membuat artikel \"" + title + "\"",
			"waktu":   getTimeAgo(p.CreatedAt),
			"tanggal": p.CreatedAt.Format("2006-01-02 15:04:05"),
			"tipe":    "artikel",
		})
	}

	var karya []models.Karya
	config.DB.Order("createdAt desc").Limit(3).Find(&karya)

	for _, k := range karya {
		aktivitas = append(aktivitas, gin.H{
			"id":      k.ID,
			"nama":    k.NamaPembuat,
			"aksi":    "mengajukan karya",
			"waktu":   getTimeAgo(k.CreatedAt),
			"tanggal": k.CreatedAt.Format("2006-01-02 15:04:05"),
			"tipe":    "karya",
		})
	}

	var alumni []models.Alumni
	config.DB.Order("created_at desc").Limit(3).Find(&alumni)

	for _, a := range alumni {
		aktivitas = append(aktivitas, gin.H{
			"id":      a.ID,
			"nama":    a.NamaLengkap,
			"aksi":    "mengajukan data alumni",
			"waktu":   getTimeAgo(a.CreatedAt),
			"tanggal": a.CreatedAt.Format("2006-01-02 15:04:05"),
			"tipe":    "alumni",
		})
	}

	// =====================
	// RESPONSE
	// =====================

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"totalSiswaAktif": totalSiswa,

			"artikel": gin.H{
				"total":    totalArtikel,
				"growth":   artikelGrowth,
				"bulanIni": artikelBulanIni,
			},

			"guruBK": gin.H{
				"total":  totalGuru,
				"growth": guruGrowth,
			},

			"karyaSiswa": gin.H{
				"total":  totalKarya,
				"growth": karyaGrowth,
			},

			"alumni": gin.H{
				"total":  totalAlumni,
				"growth": alumniGrowth,
			},

			"kuliah":    kuliah,
			"bekerja":   bekerja,
			"wirausaha": wirausaha,

			"pendingKarya":  pendingKarya,
			"pendingAlumni": pendingAlumni,

			"aktivitasTerbaru": aktivitas,
		},
	})
}

func calculateGrowth(current, previous int64) int64 {
	if previous == 0 {
		return 100
	}
	return int64(((float64(current-previous) / float64(previous)) * 100))
}

func getTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)

	if minutes < 1 {
		return "Baru saja"
	}
	if minutes < 60 {
		return fmt.Sprintf("%d menit lalu", minutes)
	}
	if hours < 24 {
		return fmt.Sprintf("%d jam lalu", hours)
	}
	if days < 30 {
		return fmt.Sprintf("%d hari lalu", days)
	}

	return t.Format("02 Jan 2006")
}
