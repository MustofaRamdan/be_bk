package controllers

import (
	"backend-bk/config"
	"backend-bk/models"
	"backend-bk/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ======================
// GET ALL KONSELING
// ======================

func GetKonseling(c *gin.Context) {
	var konseling []models.Konseling

	err := config.DB.
		Order("createdAt desc").
		Find(&konseling).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal mengambil data konseling",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    konseling,
	})
}

// ======================
// GET DETAIL KONSELING
// ======================

func GetKonselingByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var konseling models.Konseling

	err = config.DB.
		First(&konseling, id).Error

	if err != nil {
		c.JSON(404, gin.H{
			"error": "Konseling tidak ditemukan",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    konseling,
	})
}

// ======================
// PUT JAWABAN KONSELING
// ======================

func UpdateJawabanKonseling(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var body struct {
		Jawaban string `json:"jawaban"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": "Body tidak valid",
		})
		return
	}

	if strings.TrimSpace(body.Jawaban) == "" {
		c.JSON(400, gin.H{
			"error": "Jawaban tidak boleh kosong",
		})
		return
	}

	var konseling models.Konseling

	if err := config.DB.First(&konseling, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Konseling tidak ditemukan",
		})
		return
	}

	// ======================
	// UPDATE JAWABAN
	// ======================

	jawaban := strings.TrimSpace(body.Jawaban)

	konseling.Jawaban = &jawaban
	konseling.Status = "SELESAI"

	if err := config.DB.Save(&konseling).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal menyimpan jawaban",
		})
		return
	}

	// ======================
	// KIRIM EMAIL
	// ======================

	if konseling.Email != nil &&
		*konseling.Email != "" &&
		konseling.InginJawabanEmail {

		bodyEmail := fmt.Sprintf(`
			<h2>Jawaban Konseling BK</h2>

			<p>Halo %s,</p>

			<p>Jawaban dari Guru BK untuk konsultasi Anda:</p>

			<hr>

			<p><strong>Topik:</strong> %s</p>

			<p><strong>Keluhan:</strong></p>

			<p>%s</p>

			<br>

			<p><strong>Jawaban Guru BK:</strong></p>

			<p>%s</p>

			<hr>

			<p>Tetap semangat 😊</p>
		`,
			konseling.Nama,
			konseling.Topik,
			konseling.Deskripsi,
			jawaban,
		)

		err := utils.SendEmail(
			*konseling.Email,
			"Jawaban Konseling BK",
			bodyEmail,
		)

		if err != nil {
			fmt.Println("❌ Gagal mengirim email:", err)
		} else {
			fmt.Println("✅ Email berhasil dikirim")
		}
	}

	// ======================
	// RESPONSE
	// ======================

	c.JSON(200, gin.H{
		"success": true,
		"message": "Jawaban berhasil disimpan",
		"data":    konseling,
	})
}

func CreateKonseling(c *gin.Context) {
	var body struct {
		Mode              string `json:"mode"`
		Nama              string `json:"nama"`
		Kelas             string `json:"kelas"`
		Jurusan           string `json:"jurusan"`
		Email             string `json:"email"`
		Topik             string `json:"topik"`
		Deskripsi         string `json:"deskripsi"`
		InginJawabanEmail bool   `json:"inginJawabanEmail"`
	}

	// ======================
	// BIND JSON
	// ======================

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": "Body tidak valid",
		})
		return
	}

	// ======================
	// VALIDASI UMUM
	// ======================

	if body.Topik == "" || body.Deskripsi == "" {
		c.JSON(400, gin.H{
			"error": "Topik dan deskripsi wajib diisi",
		})
		return
	}

	// ======================
	// VALIDASI MODE TERDAFTAR
	// ======================

	if body.Mode == "terdaftar" {
		if body.Nama == "" ||
			body.Kelas == "" ||
			body.Jurusan == "" {

			c.JSON(400, gin.H{
				"error": "Nama, kelas, dan jurusan wajib diisi untuk mode terdaftar",
			})
			return
		}
	}

	// ======================
	// VALIDASI EMAIL
	// ======================

	if body.InginJawabanEmail && body.Email == "" {
		c.JSON(400, gin.H{
			"error": "Email wajib diisi jika ingin menerima jawaban via email",
		})
		return
	}

	// ======================
	// DEFAULT VALUE
	// ======================

	nama := body.Nama
	kelas := body.Kelas
	jurusan := body.Jurusan

	if nama == "" {
		nama = "Anonim"
	}

	if kelas == "" {
		kelas = "-"
	}

	if jurusan == "" {
		jurusan = "-"
	}

	// ======================
	// CREATE DATA
	// ======================

	konseling := models.Konseling{
		Mode:              body.Mode,
		Nama:              nama,
		Kelas:             kelas,
		Jurusan:           jurusan,
		Email:             &body.Email,
		Topik:             body.Topik,
		Deskripsi:         body.Deskripsi,
		InginJawabanEmail: body.InginJawabanEmail,
		Status:            "MENUNGGU",
	}

	// kalau email kosong -> nil
	if body.Email == "" {
		konseling.Email = nil
	}

	if err := config.DB.Create(&konseling).Error; err != nil {
		c.JSON(500, gin.H{
			"error":  "Gagal mengirim konseling",
			"detail": err.Error(),
		})
		return
	}

	// ======================
	// RESPONSE
	// ======================

	c.JSON(201, gin.H{
		"success": true,
		"message": "Konseling berhasil dikirim",
		"data":    konseling,
	})
}