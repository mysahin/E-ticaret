package Models

import "gorm.io/gorm"

type Login struct {
	gorm.Model
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
