package Models

import "gorm.io/gorm"

type Files struct {
	gorm.Model
	ProductId uint64 `json:"product_id"`
	FileName  string `json:"file_name"`
	FileUrl   string `json:"file_url"`
}
