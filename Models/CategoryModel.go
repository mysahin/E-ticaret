package Models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"category_name"`
}
