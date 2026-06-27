package controllers

import (
	"backend-bk/config"
	"backend-bk/models"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Thumbnail   string `json:"thumbnail"`
		PublishedAt string `json:"publishedAt"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Body tidak valid"})
		return
	}

	if strings.TrimSpace(body.Title) == "" {
		c.JSON(400, gin.H{"error": "Judul wajib diisi"})
		return
	}

	if body.Content == "" {
		c.JSON(400, gin.H{"error": "Konten wajib diisi"})
		return
	}

	now := time.Now()
	publishedDate := &now

	post := models.Post{
		Title:       body.Title,
		Content:     body.Content,
		Thumbnail:   &body.Thumbnail,
		Published:   true,
		PublishedAt: publishedDate,
	}

	config.DB.Create(&post)

	c.JSON(200, gin.H{
		"success": true,
		"data": post,
	})
}
func GetPosts(c *gin.Context) {
	var posts []models.Post

	result := config.DB.
		Where("published = ?", true).
		Order("createdAt desc").
		Find(&posts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil artikel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    posts,
	})
}

func GetPostByID(c *gin.Context) {
	id := c.Param("id")

	var post models.Post

	result := config.DB.First(&post, id)

	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": post,
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post

	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	var body struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Thumbnail   string `json:"thumbnail"`
		PublishedAt string `json:"publishedAt"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Body tidak valid"})
		return
	}

	if strings.TrimSpace(body.Title) == "" {
		c.JSON(400, gin.H{"error": "Judul wajib diisi"})
		return
	}

	if body.Content == "" || body.Content == "<p></p>" {
		c.JSON(400, gin.H{"error": "Konten wajib diisi"})
		return
	}

	post.Title = strings.TrimSpace(body.Title)
	post.Content = body.Content
	if post.Thumbnail != nil {
	os.Remove("." + *post.Thumbnail)
}

	if body.Thumbnail != "" {
		post.Thumbnail = &body.Thumbnail
	}

	if body.PublishedAt != "" {
		t, err := time.Parse("2006-01-02", body.PublishedAt)
		if err == nil {
			post.PublishedAt = &t
		}
	}

	config.DB.Save(&post)
	
	c.JSON(200, gin.H{
		"success": true,
		"message": "Artikel berhasil diperbarui",
		"data": post,
	})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}
	// Hapus thumbnail utama
	if post.Thumbnail != nil && *post.Thumbnail != "" {
		path := "." + *post.Thumbnail
		os.Remove(path)
	}

	// Hapus gambar yang di-upload di dalam konten editor artikel
	deleteContentImages(post.Content)

	config.DB.Delete(&post)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Artikel dihapus",
	})
}

func deleteContentImages(content string) {
	re := regexp.MustCompile(`src="([^"]*?/uploads/[^"]+)"`)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			urlPath := match[1]
			idx := strings.Index(urlPath, "/uploads/")
			if idx != -1 {
				filePath := "." + urlPath[idx:]
				os.Remove(filePath)
			}
		}
	}
}