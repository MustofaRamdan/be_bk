package controllers

import (
	"backend-bk/config"
	"backend-bk/models"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// dashAlumniStat for GROUP BY alumni status queries
type dashAlumniStat struct {
	Status string
	Total  int64
}

// growthResult for conditional COUNT queries
type growthResult struct {
	BulanIni  int64
	BulanLalu int64
}

func GetDashboard(c *gin.Context) {
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	firstDayOfLastMonth := firstDayOfMonth.AddDate(0, -1, 0)

	var wg sync.WaitGroup

	var totalSiswa int64
	var totalArtikel int64
	var totalGuru int64
	var totalKarya int64
	var alumniStats []dashAlumniStat
	var pendingKarya int64
	var pendingAlumni int64
	var artikelGrowthResult growthResult
	var guruGrowthResult growthResult
	var karyaGrowthResult growthResult
	var alumniGrowthResult growthResult
	var posts []models.Post
	var karya []models.Karya
	var alumni []models.Alumni

	wg.Add(14)

	// 1. Total Siswa
	go func() {
		defer wg.Done()
		config.DB.Model(&models.User{}).Count(&totalSiswa)
	}()

	// 2. Total Artikel
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Post{}).Where("published = ?", true).Count(&totalArtikel)
	}()

	// 3. Total Guru
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Guru{}).Count(&totalGuru)
	}()

	// 4. Total Karya
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Karya{}).Count(&totalKarya)
	}()

	// 5. Alumni Stats
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Alumni{}).
			Select("status, COUNT(*) as total").
			Where("status_pengajuan = ?", "DITERIMA").
			Group("status").
			Find(&alumniStats)
	}()

	// 6. Pending Karya
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Karya{}).Where("status = ?", "PENDING").Count(&pendingKarya)
	}()

	// 7. Pending Alumni
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Alumni{}).Where("status_pengajuan = ?", "PENDING").Count(&pendingAlumni)
	}()

	// 8. Growth Artikel
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Post{}).
			Select(
				"COALESCE(SUM(CASE WHEN createdAt >= ? THEN 1 ELSE 0 END), 0) as bulan_ini, "+
					"COALESCE(SUM(CASE WHEN createdAt >= ? AND createdAt < ? THEN 1 ELSE 0 END), 0) as bulan_lalu",
				firstDayOfMonth, firstDayOfLastMonth, firstDayOfMonth,
			).
			Where("published = ?", true).
			Find(&artikelGrowthResult)
	}()

	// 9. Growth Guru
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Guru{}).
			Select(
				"COALESCE(SUM(CASE WHEN createdAt >= ? THEN 1 ELSE 0 END), 0) as bulan_ini, "+
					"COALESCE(SUM(CASE WHEN createdAt >= ? AND createdAt < ? THEN 1 ELSE 0 END), 0) as bulan_lalu",
				firstDayOfMonth, firstDayOfLastMonth, firstDayOfMonth,
			).
			Find(&guruGrowthResult)
	}()

	// 10. Growth Karya
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Karya{}).
			Select(
				"COALESCE(SUM(CASE WHEN createdAt >= ? THEN 1 ELSE 0 END), 0) as bulan_ini, "+
					"COALESCE(SUM(CASE WHEN createdAt >= ? AND createdAt < ? THEN 1 ELSE 0 END), 0) as bulan_lalu",
				firstDayOfMonth, firstDayOfLastMonth, firstDayOfMonth,
			).
			Find(&karyaGrowthResult)
	}()

	// 11. Growth Alumni
	go func() {
		defer wg.Done()
		config.DB.Model(&models.Alumni{}).
			Select(
				"COALESCE(SUM(CASE WHEN created_at >= ? THEN 1 ELSE 0 END), 0) as bulan_ini, "+
					"COALESCE(SUM(CASE WHEN created_at >= ? AND created_at < ? THEN 1 ELSE 0 END), 0) as bulan_lalu",
				firstDayOfMonth, firstDayOfLastMonth, firstDayOfMonth,
			).
			Where("status_pengajuan = ?", "DITERIMA").
			Find(&alumniGrowthResult)
	}()

	// 12. Recent Posts
	go func() {
		defer wg.Done()
		config.DB.Order("createdAt desc").Limit(3).Find(&posts)
	}()

	// 13. Recent Karya
	go func() {
		defer wg.Done()
		config.DB.Order("createdAt desc").Limit(3).Find(&karya)
	}()

	// 14. Recent Alumni
	go func() {
		defer wg.Done()
		config.DB.Order("created_at desc").Limit(3).Find(&alumni)
	}()

	wg.Wait()

	var totalAlumni, kuliah, bekerja, wirausaha int64
	for _, s := range alumniStats {
		totalAlumni += s.Total
		switch s.Status {
		case "KULIAH":
			kuliah = s.Total
		case "BEKERJA":
			bekerja = s.Total
		case "WIRAUSAHA":
			wirausaha = s.Total
		}
	}

	artikelGrowth := calculateGrowth(artikelGrowthResult.BulanIni, artikelGrowthResult.BulanLalu)
	guruGrowth := calculateGrowth(guruGrowthResult.BulanIni, guruGrowthResult.BulanLalu)
	karyaGrowth := calculateGrowth(karyaGrowthResult.BulanIni, karyaGrowthResult.BulanLalu)
	alumniGrowth := calculateGrowth(alumniGrowthResult.BulanIni, alumniGrowthResult.BulanLalu)

	var aktivitas []gin.H

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
				"bulanIni": artikelGrowthResult.BulanIni,
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
