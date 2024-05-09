package Models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	TypeId           int     `json:"type_id"`
	ProductName      string  `json:"product_name"`
	ProductPrice     int     `json:"product_price"`
	ProductTitle     string  `json:"product_title"`
	ProductStatement string  `json:"product_statement"`
	SellerUserName   string  `json:"seller_username"`
	ProductCount     int     `json:"product_count"`
	Archived         bool    `json:"archived"`
	ProductRating    float64 `json:"product_rating"`
	NumberOfRatings  int     `json:"number_of_ratings"`
}
