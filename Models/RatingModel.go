package Models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	ProductId uint64 `json:"product_id"`
	Username  string `json:"username"`
	Rating    int    `json:"rating"`
}
