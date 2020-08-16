package entity

import "time"

// Item struct is main struct for table "m_barang"
type Item struct {
	ID           int       `json:"id" db:"id"`
	CategoryID   int       `json:"category_id" db:"kategori_id"`
	BrandID      int       `json:"brand_id" db:"merk_id"`
	Code         string    `json:"code" db:"kode"`
	Name         string    `json:"name" db:"nama"`
	Length       float64   `json:"length" db:"panjang"`
	Width        float64   `json:"width" db:"lebar"`
	InitialStock float64   `json:"initial_stock" db:"initial_stok"`
	InitialArea  float64   `json:"initial_area" db:"initial_m2"`
	Remark       string    `json:"remark" db:"keterangan"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
	UpdatedBy    string    `json:"updated_by" db:"updated_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
