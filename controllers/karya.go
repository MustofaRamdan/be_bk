	package controllers

	import (
		"backend-bk/config"
		"backend-bk/models"
		"backend-bk/utils"
		"fmt"
		"os"
		"strings"

		"github.com/gin-gonic/gin"
	)

	// ======================
	// CREATE KARYA
	// ======================

	func CreateKarya(c *gin.Context) {
		var body struct {
			Judul       string `json:"judul"`
			Deskripsi   string `json:"deskripsi"`
			Link        string `json:"link"`
			NamaPembuat string `json:"namaPembuat"`
			Kelas       string `json:"kelas"`
			Jurusan     string `json:"jurusan"`
			Email       string `json:"email"`
			Thumbnail   string `json:"thumbnail"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{
				"error": "Body tidak valid",
			})
			return
		}

		if strings.TrimSpace(body.Judul) == "" {
			c.JSON(400, gin.H{
				"error": "Judul karya wajib diisi",
			})
			return
		}

		karya := models.Karya{
			Judul:       body.Judul,
			Deskripsi:   body.Deskripsi,
			Link:        &body.Link,
			NamaPembuat: body.NamaPembuat,
			Kelas:       body.Kelas,
			Jurusan:     body.Jurusan,
			Thumbnail:   &body.Thumbnail,
			Status:      "PENDING",
		}

		// email optional
		if body.Email != "" {
			karya.Email = &body.Email
		}

		if err := config.DB.Create(&karya).Error; err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"data":    karya,
		})
	}

	// ======================
	// GET ALL KARYA
	// ======================

	func GetKarya(c *gin.Context) {
		status := c.Query("status")

		var karya []models.Karya

		query := config.DB.Order("createdAt desc")

		if status != "" {
			query = query.Where("status = ?", status)
		}

		if err := query.Find(&karya).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "Gagal mengambil data",
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    karya,
		})
	}

	// ======================
	// GET KARYA BY ID
	// ======================

	func GetKaryaByID(c *gin.Context) {
		id := c.Param("id")

		var karya models.Karya

		if err := config.DB.First(&karya, id).Error; err != nil {
			c.JSON(404, gin.H{
				"error": "Karya tidak ditemukan",
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    karya,
		})
	}

	// ======================
	// UPDATE KARYA
	// ======================

// ======================
// UPDATE KARYA
// ======================

func UpdateKarya(c *gin.Context) {
	id := c.Param("id")

	var karya models.Karya

	if err := config.DB.First(&karya, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Karya tidak ditemukan",
		})
		return
	}

	oldStatus := karya.Status

	var body map[string]interface{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": "Body tidak valid",
		})
		return
	}

	// ======================
	// VALIDASI STATUS
	// ======================

	if status, ok := body["status"]; ok {
		valid := map[string]bool{
			"PENDING":  true,
			"DITERIMA": true,
			"DITOLAK":  true,
		}

		if !valid[status.(string)] {
			c.JSON(400, gin.H{
				"error": "Status tidak valid",
			})
			return
		}
	}

	// ======================
	// VALIDASI ALASAN (wajib kalau DITOLAK)
	// ======================

	if status, ok := body["status"]; ok && status == "DITOLAK" {
		keterangan, hasKeterangan := body["keterangan"]
		if !hasKeterangan || keterangan == nil || strings.TrimSpace(keterangan.(string)) == "" {
			c.JSON(400, gin.H{
				"error": "Alasan penolakan wajib diisi",
			})
			return
		}
	}

	// ======================
	// HAPUS THUMBNAIL LAMA
	// ======================

	if newThumb, exists := body["thumbnail"]; exists {
		if karya.Thumbnail != nil && *karya.Thumbnail != "" {
			if newPath, ok := newThumb.(string); ok {
				if newPath != *karya.Thumbnail {
					os.Remove("." + *karya.Thumbnail)
				}
			}
		}
	}

	// ======================
	// UPDATE DB
	// ======================

	if err := config.DB.Model(&karya).Updates(body).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal mengupdate karya",
		})
		return
	}

	// ambil data terbaru
	config.DB.First(&karya, id)

	// ======================
	// KIRIM EMAIL NOTIFIKASI
	// ======================

	if karya.Email != nil && *karya.Email != "" && oldStatus != karya.Status {
		var subject string
		var bodyEmail string

		// DITERIMA
		if karya.Status == "DITERIMA" {
			subject = "Karya Anda Diterima 🎉"

			bodyEmail = fmt.Sprintf(`
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
					<div style="background: #6b7c4e; padding: 20px; text-align: center;">
						<h1 style="color: white; margin: 0;">🎉 Karya Diterima!</h1>
					</div>
					<div style="padding: 20px; background: #f9f9f9;">
						<p>Halo <strong>%s</strong>,</p>
						<p>Selamat! Karya Anda telah <strong style="color: #22c55e;">DITERIMA</strong>.</p>
						
						<div style="background: white; padding: 16px; border-radius: 8px; margin: 16px 0;">
							<h3 style="margin: 0 0 8px 0; color: #333;">%s</h3>
							<p style="color: #666; margin: 0;">%s</p>
						</div>

						<p>Karya Anda sekarang ditampilkan di halaman Karya Siswa.</p>
						<p>Terima kasih telah berkontribusi! 😊</p>
					</div>
				</div>
			`, karya.NamaPembuat, karya.Judul, karya.Deskripsi)
		}

		// DITOLAK
		if karya.Status == "DITOLAK" {
			subject = "Karya Anda Perlu Perbaikan"

			// Ambil alasan penolakan
			alasan := "Tidak ada keterangan"
			if karya.Keterangan != nil && *karya.Keterangan != "" {
				alasan = *karya.Keterangan
			}

			bodyEmail = fmt.Sprintf(`
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
					<div style="background: #dc2626; padding: 20px; text-align: center;">
						<h1 style="color: white; margin: 0;">Karya Perlu Perbaikan</h1>
					</div>
					<div style="padding: 20px; background: #f9f9f9;">
						<p>Halo <strong>%s</strong>,</p>
						<p>Mohon maaf, karya Anda belum dapat kami tampilkan.</p>
						
						<div style="background: white; padding: 16px; border-radius: 8px; margin: 16px 0;">
							<h3 style="margin: 0 0 8px 0; color: #333;">%s</h3>
							<p style="color: #666; margin: 0;">%s</p>
						</div>

						<div style="background: #fef2f2; border-left: 4px solid #dc2626; padding: 16px; margin: 16px 0;">
							<h4 style="color: #dc2626; margin: 0 0 8px 0;">Alasan:</h4>
							<p style="color: #991b1b; margin: 0;">%s</p>
						</div>

						<p>Silakan perbaiki sesuai arahan di atas dan kirim ulang karya Anda.</p>
						<p>Tetap semangat berkarya! 💪</p>
					</div>
				</div>
			`, karya.NamaPembuat, karya.Judul, karya.Deskripsi, alasan)
		}

		// Kirim email
		if bodyEmail != "" {
			err := utils.SendEmail(*karya.Email, subject, bodyEmail)
			if err != nil {
				fmt.Println("❌ Gagal kirim email:", err)
			} else {
				fmt.Println("✅ Email berhasil dikirim ke:", *karya.Email)
			}
		}
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    karya,
	})
}

	// ======================
	// DELETE KARYA
	// ======================

	func DeleteKarya(c *gin.Context) {
		id := c.Param("id")

		var karya models.Karya

		if err := config.DB.First(&karya, id).Error; err != nil {
			c.JSON(404, gin.H{
				"error": "Karya tidak ditemukan",
			})
			return
		}

		// hapus thumbnail
		if karya.Thumbnail != nil &&
			*karya.Thumbnail != "" {

			os.Remove("." + *karya.Thumbnail)
		}

		config.DB.Delete(&karya)

		c.JSON(200, gin.H{
			"success": true,
			"message": "Karya dihapus",
		})
	}