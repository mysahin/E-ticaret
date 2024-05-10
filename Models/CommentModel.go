package Models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ProductId uint64 `json:"product_id"`
	Username  string `json:"username"`
	Comment   string `json:"comment"`
}
