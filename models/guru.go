package models

import "time"

type Guru struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	Nama      string    `json:"nama" gorm:"column:nama"`
	Jabatan   string    `json:"jabatan" gorm:"column:jabatan"`
	Kelas     string    `json:"kelas" gorm:"column:kelas"`
	Foto      *string   `json:"foto" gorm:"column:foto"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt"`
}

func (Guru) TableName() string {
	return "guru"
}