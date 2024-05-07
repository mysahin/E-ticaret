package Controllers

import (
	database "ETicaret/Database"
	"ETicaret/Helpers"
	"ETicaret/Models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Product struct{}

func (product Product) AddProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		db := database.DB.Db
		product := new(Models.Product)
		if err := c.BodyParser(&product); err != nil {
			return err
		}
		userName := getUserName(c)
		newProduct := Models.Product{
			TypeId:           product.TypeId,
			ProductName:      product.ProductName,
			ProductPrice:     product.ProductPrice,
			ProductStatement: product.ProductStatement,
			ProductTitle:     product.ProductTitle,
			SellerUserName:   userName,
			ProductCount:     product.ProductCount,
		}
		if err := db.Create(&newProduct).Error; err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"message": "Yeni ürün başarıyla eklendi.",
			"product": newProduct,
		})
	}

	return c.JSON(fiber.Map{
		"message": "hata",
	})
}

func (product Product) ViewMyProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		db := database.DB.Db
		var products []Models.Product
		username := getUserName(c)
		if err := db.Find(&products, "seller_user_name=?", username).Error; err != nil {
		}
		fmt.Println(getID(c))
		return c.JSON(fiber.Map{
			"Ürünleriniz": products,
		})
	}

	return c.JSON(fiber.Map{
		"Message": "Önce giriş yapmalısınız!!!",
	})

}
