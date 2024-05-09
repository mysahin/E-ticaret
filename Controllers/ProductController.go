package Controllers

import (
	database "ETicaret/Database"
	"ETicaret/Helpers"
	"ETicaret/Models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type Product struct{}

func (product Product) AddProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		db := database.DB.Db
		addedProduct := new(Models.Product)
		if err := c.BodyParser(&product); err != nil {
			return err
		}
		userName := getUserName(c)
		newProduct := Models.Product{
			TypeId:           addedProduct.TypeId,
			ProductName:      addedProduct.ProductName,
			ProductPrice:     addedProduct.ProductPrice,
			ProductStatement: addedProduct.ProductStatement,
			ProductTitle:     addedProduct.ProductTitle,
			SellerUserName:   userName,
			ProductCount:     addedProduct.ProductCount,
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
		if err := db.Find(&products, "seller_user_name=? AND archived=?", username, "0").Error; err != nil {
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
func (product Product) ViewProductById(c *fiber.Ctx) error {
	db := database.DB.Db
	var products Models.Product
	productId := c.Params("id")
	if err := db.First(&products, "id=? AND archived=?", productId, "0").Error; err != nil {
		return c.JSON(fiber.Map{
			"error": err,
		})
	}
	return c.JSON(fiber.Map{
		"Ürün": products,
	})
}

func (product Product) ViewProductsByType(c *fiber.Ctx) error {
	db := database.DB.Db
	var products []Models.Product
	productType := c.Params("type")
	if err := db.Find(&products, "type_id=? AND archived=?", productType, "0").Error; err != nil {
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
		if err := db.Find(&products, "type_id=? AND archived=?", x.ID, "0").Error; err != nil {
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

func (product Product) DeleteProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		db := database.DB.Db
		productId := c.Params("id")
		var deleteProduct Models.Product

		if err := db.First(&deleteProduct, "id = ?", productId).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Gönderi bulunamadı.",
			})
		}
		username := getUserName(c)
		if deleteProduct.SellerUserName == username {
			if err := db.Delete(&deleteProduct).Error; err != nil {
				return err
			}

			return c.JSON(fiber.Map{
				"message": "Ürün başarıyla silindi.",
			})
		}
		return c.JSON(fiber.Map{
			"message": "Bu ürünü silme yetkiniz yok!",
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Lütfen önce giriş yapınız.",
	})
}

func (product Product) ArchiveProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		username := getUserName(c)

		db := database.DB.Db
		postID := c.Params("id")
		var archivedProduct Models.Product

		if err := db.First(&archivedProduct, postID).Error; err != nil {
			return err
		}
		if archivedProduct.SellerUserName == username {
			if err := db.First(&archivedProduct).Where("id=?", postID).Error; err != nil {
				return err
			}
			if archivedProduct.Archived == false {
				if err := db.Model(&archivedProduct).Where("id=?", postID).Update("archived", "1").Error; err != nil {
					return err
				}
				return c.JSON("Başarıyla arşivlendi.")
			} else {
				if err := db.Model(&archivedProduct).Where("id=?", postID).Update("archived", "0").Error; err != nil {
					return err
				}
				return c.JSON("Başarıyla arşivden çıkarıldı.")
			}

		}
		return c.JSON(fiber.Map{
			"Warning": "Bu ürün için yetkiniz yok!",
		})
	}

	return c.JSON("lütfen giriş yapınız!!!")
}

func (product Product) EditProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		db := database.DB.Db
		editedProduct := new(Models.Product)
		if err := c.BodyParser(&editedProduct); err != nil {
			return err
		}
		if err := db.Model(&editedProduct).Where("id=?", editedProduct.ID).Update("product_statement", editedProduct.ProductStatement).Error; err != nil {
			return err
		}
		if err := db.Model(&editedProduct).Where("id=?", editedProduct.ID).Update("product_title", editedProduct.ProductTitle).Error; err != nil {
			return err
		}
		if err := db.Model(&editedProduct).Where("id=?", editedProduct.ID).Update("product_name", editedProduct.ProductName).Error; err != nil {
			return err
		}
		if err := db.Model(&editedProduct).Where("id=?", editedProduct.ID).Update("product_price", editedProduct.ProductPrice).Error; err != nil {
			return err
		}
		if err := db.Model(&editedProduct).Where("id=?", editedProduct.ID).Update("product_count", editedProduct.ProductCount).Error; err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"message": "Ürün başarıyla güncellendi.",
			"product": editedProduct,
		})
	}

	return c.JSON(fiber.Map{
		"Warning": "Lütfen giriş yapınız!",
	})
}

func (product Product) RateProduct(c *fiber.Ctx) error {
	productID := c.Params("productID")
	rating, err := strconv.Atoi(c.Params("rating"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz puanlama formatı",
		})
	}

	uIntID, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		return err
	}

	db := database.DB.Db
	username := getUserName(c)

	// Rating tablosunda kullanıcının daha önce bu ürüne puan verip vermediğini kontrol et
	var existingRating Models.Rating
	if err := db.Where("username = ? AND product_id = ?", username, uIntID).First(&existingRating).Error; err == nil {
		return c.JSON(fiber.Map{
			"message": "Bu ürüne zaten puan verdiniz.",
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Yeni bir rating oluştur
	newRating := Models.Rating{
		ProductId: uIntID,
		Username:  username,
		Rating:    rating,
	}

	if err := db.Create(&newRating).Error; err != nil {
		return err
	}

	// Rating tablosunda, verilen ürün ID'sine sahip rating'lerin ortalamasını al
	var averageRating float64
	if err := db.Model(&Models.Rating{}).Where("product_id = ?", productID).Select("AVG(rating) as average_rating").Scan(&averageRating).Error; err != nil {
		return err
	}

	// Ürünün rating'ini güncelle
	if err := db.Model(&Models.Product{}).Where("id = ?", productID).Update("product_rating", averageRating).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Ürüne başarıyla puan verdiniz.",
	})
}
