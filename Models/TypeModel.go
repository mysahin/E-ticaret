package Models

import "gorm.io/gorm"

type Type struct {
	gorm.Model
	Name       string `json:"type_name"`
	CategoryId int    `json:"category_id"`
}
