package models

import "time"

type Konseling struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Mode              string    `json:"mode"`
	Nama              string    `json:"nama"`
	Kelas             string    `json:"kelas"`
	Jurusan           string    `json:"jurusan"`
	Email             *string   `json:"email"`
	Topik             string    `json:"topik"`
	Deskripsi         string    `gorm:"type:text" json:"deskripsi"`
	InginJawabanEmail bool      `json:"inginJawabanEmail"`
	Status            string    `json:"status"`
	Jawaban           *string   `gorm:"type:text" json:"jawaban"`

	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}