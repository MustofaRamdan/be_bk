package models

import "time"

type Karya struct {
	ID          uint      `json:"id" gorm:"primaryKey;column:id"`
	Judul       string    `json:"judul" gorm:"column:judul"`
	Deskripsi   string    `json:"deskripsi" gorm:"column:deskripsi"`
	Link        *string   `json:"link" gorm:"column:link"`
	NamaPembuat string    `json:"namaPembuat" gorm:"column:namaPembuat"`
	Kelas       string    `json:"kelas" gorm:"column:kelas"`
	Email  		*string 	  `json:"email"`
	Jurusan     string    `json:"jurusan" gorm:"column:jurusan"`
	Thumbnail   *string   `json:"thumbnail" gorm:"column:thumbnail"`
	Status      string    `json:"status" gorm:"column:status;index"`
	Keterangan  *string   `json:"keterangan" gorm:"column:keterangan"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:createdAt;index"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updatedAt"`
}

func (Karya) TableName() string {
	return "karya"
}