package controllers

import (
	"backend-bk/config"
	"backend-bk/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateAlumni(c *gin.Context) {
	var body struct {
		NamaLengkap       string `json:"namaLengkap"`
		TahunLulus        string `json:"tahunLulus"`
		Status            string `json:"status"`

		NamaKampus        string `json:"namaKampus"`
		ProgramStudi      string `json:"programStudi"`
		TahunMasukKuliah  string `json:"tahunMasukKuliah"`

		NamaPerusahaan    string `json:"namaPerusahaan"`
		TahunMasukKerja   string `json:"tahunMasukKerja"`

		NamaUsaha         string `json:"namaUsaha"`
		TahunAwalUsaha    string `json:"tahunAwalUsaha"`

		BuktiPendukung    string `json:"buktiPendukung"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body tidak valid"})
		return
	}

	// Validasi umum
	if strings.TrimSpace(body.NamaLengkap) == "" {
		c.JSON(400, gin.H{"error": "Nama lengkap wajib diisi"})
		return
	}

	if strings.TrimSpace(body.TahunLulus) == "" {
		c.JSON(400, gin.H{"error": "Tahun lulus wajib diisi"})
		return
	}

	if body.Status == "" {
		c.JSON(400, gin.H{"error": "Status wajib dipilih"})
		return
	}

	// Validasi berdasarkan status
	switch body.Status {
	case "KULIAH":
		if strings.TrimSpace(body.NamaKampus) == "" {
			c.JSON(400, gin.H{"error": "Nama kampus wajib diisi"})
			return
		}
		if strings.TrimSpace(body.ProgramStudi) == "" {
			c.JSON(400, gin.H{"error": "Program studi wajib diisi"})
			return
		}
		if strings.TrimSpace(body.TahunMasukKuliah) == "" {
			c.JSON(400, gin.H{"error": "Tahun masuk kuliah wajib diisi"})
			return
		}

	case "BEKERJA":
		if strings.TrimSpace(body.NamaPerusahaan) == "" {
			c.JSON(400, gin.H{"error": "Nama perusahaan wajib diisi"})
			return
		}
		if strings.TrimSpace(body.TahunMasukKerja) == "" {
			c.JSON(400, gin.H{"error": "Tahun masuk kerja wajib diisi"})
			return
		}

	case "WIRAUSAHA":
		if strings.TrimSpace(body.NamaUsaha) == "" {
			c.JSON(400, gin.H{"error": "Nama usaha wajib diisi"})
			return
		}
		if strings.TrimSpace(body.TahunAwalUsaha) == "" {
			c.JSON(400, gin.H{"error": "Tahun awal usaha wajib diisi"})
			return
		}

	default:
		c.JSON(400, gin.H{"error": "Status tidak valid"})
		return
	}

	// helper pointer
	strPtr := func(s string) *string {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		return &s
	}

	alumni := models.Alumni{
		NamaLengkap:       strings.TrimSpace(body.NamaLengkap),
		TahunLulus:        strings.TrimSpace(body.TahunLulus),
		Status:            models.StatusAlumni(body.Status),

		NamaKampus:        strPtr(body.NamaKampus),
		ProgramStudi:      strPtr(body.ProgramStudi),
		TahunMasukKuliah:  strPtr(body.TahunMasukKuliah),

		NamaPerusahaan:    strPtr(body.NamaPerusahaan),
		TahunMasukKerja:   strPtr(body.TahunMasukKerja),

		NamaUsaha:         strPtr(body.NamaUsaha),
		TahunAwalUsaha:    strPtr(body.TahunAwalUsaha),

		BuktiPendukung:    strPtr(body.BuktiPendukung),
		StatusPengajuan:   models.PENDING,
	}

	if err := config.DB.Create(&alumni).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"data":    alumni,
	})
}

func GetAlumni(c *gin.Context) {
	var alumni []models.Alumni

	query := config.DB.Order("created_at desc")

	// filter status pengajuan
	statusPengajuan := c.Query("statusPengajuan")

	if statusPengajuan != "" {
		query = query.Where("status_pengajuan = ?", statusPengajuan)
	}

	// filter status alumni
	status := c.Query("status")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&alumni).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal mengambil data",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    alumni,
	})
}

func GetAlumniByID(c *gin.Context) {
	id := c.Param("id")

	var alumni models.Alumni

	if err := config.DB.First(&alumni, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Alumni tidak ditemukan"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    alumni,
	})
}

func UpdateAlumniStatus(c *gin.Context) {
	id := c.Param("id")

	var alumni models.Alumni

	if err := config.DB.First(&alumni, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Alumni tidak ditemukan"})
		return
	}

	var body struct {
		StatusPengajuan string `json:"statusPengajuan"`
		Keterangan      string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Body tidak valid"})
		return
	}

	// Validasi status
	validStatus := map[string]bool{
		"PENDING":  true,
		"DITERIMA": true,
		"DITOLAK":  true,
	}

	if !validStatus[body.StatusPengajuan] {
		c.JSON(400, gin.H{"error": "Status pengajuan tidak valid"})
		return
	}

	// Update field
	alumni.StatusPengajuan = models.StatusPengajuan(body.StatusPengajuan)

	if body.Keterangan != "" {
		alumni.Keterangan = &body.Keterangan
	} else {
		alumni.Keterangan = nil
	}

	if err := config.DB.Save(&alumni).Error; err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengupdate alumni"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    alumni,
	})
}


