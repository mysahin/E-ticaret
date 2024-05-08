package Models

type Cart struct {
	Product      Product
	ProductName  string `json:"product_name"`
	Quantity     string `json:"quantity"`
	ProductId    uint   `json:"product_id"`
	ProductPrice int    `json:"product_price"`
}
