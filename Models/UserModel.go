package Models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	PhoneNumber int    `json:"phone_number"`
	Adress      string `json:"adress"`
	TcNo        int    `json:"tc_no"`
}
