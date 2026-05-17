package models

import "time"

type StatusAlumni string
type StatusPengajuan string

const (
	KULIAH    StatusAlumni = "KULIAH"
	BEKERJA   StatusAlumni = "BEKERJA"
	WIRAUSAHA StatusAlumni = "WIRAUSAHA"

	PENDING  StatusPengajuan = "PENDING"
	DITERIMA StatusPengajuan = "DITERIMA"
	DITOLAK  StatusPengajuan = "DITOLAK"
)

type Alumni struct {
	ID              uint             `json:"id" gorm:"primaryKey"`
	NamaLengkap     string           `json:"namaLengkap"`
	TahunLulus      string           `json:"tahunLulus"`
	Status          StatusAlumni     `json:"status"`

	// Kuliah
	NamaKampus       *string `json:"namaKampus"`
	ProgramStudi     *string `json:"programStudi"`
	TahunMasukKuliah *string `json:"tahunMasukKuliah"`

	// Bekerja
	NamaPerusahaan  *string `json:"namaPerusahaan"`
	TahunMasukKerja *string `json:"tahunMasukKerja"`

	// Wirausaha
	NamaUsaha      *string `json:"namaUsaha"`
	TahunAwalUsaha *string `json:"tahunAwalUsaha"`

	BuktiPendukung  *string          `json:"buktiPendukung"`
	StatusPengajuan StatusPengajuan  `json:"statusPengajuan" gorm:"default:PENDING"`
	Keterangan *string `json:"keterangan" gorm:"column:keterangan"`

	CreatedAt time.Time
	UpdatedAt time.Time
}