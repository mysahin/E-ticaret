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
		"Warning": "Önce giriş yapmalısınız!",
	})

}

func (product Product) ViewProductsByType(c *fiber.Ctx) error {
	db := database.DB.Db
	var products []Models.Product
	productType := c.Params("type")
	if err := db.Find(&products, "type_id=?", productType).Error; err != nil {
		return c.JSON(fiber.Map{
			"error": err,
		})
	}

	return c.JSON(fiber.Map{
		"Ürünler": products,
	})
}

func (product Product) ViewProductsByCategory(c *fiber.Ctx) error {
	db := database.DB.Db
	var categoryProducts []Models.Product
	var types []Models.Type
	productCategory := c.Params("category")
	if err := db.Find(&types, "category_id=?", productCategory).Error; err != nil {
		return c.JSON(fiber.Map{
			"error": err,
		})
	}

	for _, x := range types {
		var products []Models.Product
		if err := db.Find(&products, "type_id=?", x.ID).Error; err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		categoryProducts = append(categoryProducts, products...)
	}

	return c.JSON(fiber.Map{
		"Ürünler": categoryProducts,
	})
}
