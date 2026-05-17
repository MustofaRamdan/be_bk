package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(400, gin.H{"error": "File wajib dipilih"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)

	path := "./uploads/" + filename

	c.SaveUploadedFile(file, path)

	c.JSON(200, gin.H{
		"url": "/uploads/" + filename,
	})
}
