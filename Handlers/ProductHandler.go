package Handlers

import (
	database "ETicaret/Database"
	"ETicaret/Helpers"
	"ETicaret/Models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type Product struct{}

func (product Product) AddProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		db := database.DB.Db
		addedProduct := new(Models.Product)
		if err := c.BodyParser(&addedProduct); err != nil {
			return err
		}
		userName := Helpers.GetUserName(c)
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

		page, err := strconv.Atoi(c.Query("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(c.Query("pageSize", "5"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize

		var products []Models.Product
		username := Helpers.GetUserName(c)
		if err := db.Where("seller_user_name=?", username).Find(&products).Error; err != nil {
			return err
		}

		totalRecords := len(products)

		totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
		currentPage := page
		if page > totalPages {
			err := c.JSON(fiber.Map{
				"Error": "Sayfa bulunamadı.",
			})

			return err

		} else {
			nextPage := currentPage + 1
			if nextPage > totalPages {
				nextPage = totalPages
			}
			prevPage := currentPage - 1
			if prevPage < 1 {
				prevPage = 1
			}

			return c.JSON(fiber.Map{
				"totalPages":  totalPages,
				"currentPage": currentPage,
				"nextPage":    nextPage,
				"prevPage":    prevPage,
				"products":    products[offset:min(offset+pageSize, totalRecords)],
			})
		}
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

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var products []Models.Product
	productType := c.Params("type")
	if err := db.Limit(pageSize).Offset(offset).Where("type_id = ? AND archived = ?", productType, "0").Find(&products).Error; err != nil {
		return err
	}

	var totalRecords int64
	if err := db.Model(&Models.Product{}).Where("type_id = ? AND archived = ?", productType, "0").Count(&totalRecords).Error; err != nil {
		return err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	currentPage := page

	// Geçersiz sayfa kontrolü
	if currentPage > totalPages {
		return c.JSON(fiber.Map{
			"Error": "Sayfa bulunamadı.",
		})
	}

	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}
	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}

	return c.JSON(fiber.Map{
		"totalPages":  totalPages,
		"currentPage": currentPage,
		"nextPage":    nextPage,
		"prevPage":    prevPage,
		"products":    products,
	})
}

func (product Product) ViewProductsByCategory(c *fiber.Ctx) error {
	db := database.DB.Db

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var categoryProducts []Models.Product
	var types []Models.Type

	productCategory := c.Params("category")

	if err := db.Where("category_id=?", productCategory).Find(&types).Error; err != nil {
		return err
	}

	for _, x := range types {
		var products []Models.Product
		if err := db.Where("type_id=? AND archived=?", x.ID, "0").Find(&products).Error; err != nil {
			return err
		}
		categoryProducts = append(categoryProducts, products...)
	}

	totalRecords := len(categoryProducts)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	currentPage := page

	// Geçersiz sayfa kontrolü
	if currentPage > totalPages {
		return c.JSON(fiber.Map{
			"Error": "Sayfa bulunamadı.",
		})
	}

	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}
	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}

	return c.JSON(fiber.Map{
		"totalPages":  totalPages,
		"currentPage": currentPage,
		"nextPage":    nextPage,
		"prevPage":    prevPage,
		"products":    categoryProducts[offset:Min(offset+pageSize, totalRecords)],
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
		username := Helpers.GetUserName(c)

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
		username := Helpers.GetUserName(c)

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
		username := Helpers.GetUserName(c)
		editedProduct := new(Models.Product)
		checkProduct := new(Models.Product)
		if err := c.BodyParser(&editedProduct); err != nil {
			return err
		}
		if err := db.First(&checkProduct).Error; err != nil {
			return err
		}
		if checkProduct.SellerUserName == username {
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
			"Warning": "Ürünün satıcısı siz değilsiniz!",
		})
	}

	return c.JSON(fiber.Map{
		"Warning": "Lütfen giriş yapınız!",
	})
}

func (product Product) RateProduct(c *fiber.Ctx) error {
	islogin := Helpers.IsLogin(c)
	if islogin {
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
		username := Helpers.GetUserName(c)

		var existingRating Models.Rating
		if err := db.Where("username = ? AND product_id = ?", username, uIntID).First(&existingRating).Error; err == nil {
			return c.JSON(fiber.Map{
				"message": "Bu ürüne zaten puan verdiniz.",
			})
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		newRating := Models.Rating{
			ProductId: uIntID,
			Username:  username,
			Rating:    rating,
		}

		if err := db.Create(&newRating).Error; err != nil {
			return err
		}

		var averageRating float64
		if err := db.Model(&Models.Rating{}).Where("product_id = ?", productID).Select("AVG(rating) as average_rating").Scan(&averageRating).Error; err != nil {
			return err
		}

		if err := db.Model(&Models.Product{}).Where("id = ?", productID).Update("product_rating", averageRating).Error; err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"message": "Ürüne başarıyla puan verdiniz.",
		})
	}
	return c.JSON(fiber.Map{
		"error": "Lütfen önce giriş yapınız.",
	})
}

func Search(c *fiber.Ctx) error {
	searchTerm := c.Query("search")

	if len(searchTerm) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Arama kelimesi en az 3 harf içermelidir.",
		})
	}

	searchTerm = "%" + searchTerm + "%"

	var categories []Models.Category
	db := database.DB.Db
	if err := db.Where("name ILIKE ?", searchTerm).Find(&categories).Error; err != nil {
		categories = nil
	}

	var types []Models.Type
	if err := db.Where("name ILIKE ?", searchTerm).Find(&types).Error; err != nil {
		types = nil
	}

	var products []Models.Product
	if err := db.Where("product_title ILIKE ?", searchTerm).Find(&products).Error; err != nil {
		products = nil
	}

	return c.JSON(fiber.Map{
		"categories": categories,
		"types":      types,
		"products":   products,
	})
}

func (product Product) CommentProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		productID := c.Params("productID")
		comment := new(Models.Comment)
		if err := c.BodyParser(comment); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		uIntID, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			return err
		}

		db := database.DB.Db
		username := Helpers.GetUserName(c)

		var existingComment Models.Comment
		if err := db.Where("username = ? AND product_id = ?", username, uIntID).First(&existingComment).Error; err == nil {
			return c.JSON(fiber.Map{
				"message": "Bu ürüne zaten yorum yaptınız.",
			})
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		newComment := Models.Comment{
			ProductId: uIntID,
			Username:  username,
			Comment:   comment.Comment,
		}

		if err := db.Create(&newComment).Error; err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"message": "Ürüne başarıyla yorum yaptınız.",
		})
	}
	return c.JSON(fiber.Map{
		"error": "Önce giriş yapınız.",
	})
}

func ViewProductComments(c *fiber.Ctx) error {
	db := database.DB.Db

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var comments []Models.Comment
	productId := c.Params("productId")
	if err := db.Limit(pageSize).Offset(offset).Where("product_id = ?", productId).Find(&comments).Error; err != nil {
		return err
	}

	var totalRecords int64
	if err := db.Model(&Models.Comment{}).Where("product_id = ? ", productId).Count(&totalRecords).Error; err != nil {
		return err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	currentPage := page

	// Geçersiz sayfa kontrolü
	if currentPage > totalPages {
		return c.JSON(fiber.Map{
			"Error": "Sayfa bulunamadı.",
		})
	}

	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}
	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}

	return c.JSON(fiber.Map{
		"totalPages":  totalPages,
		"currentPage": currentPage,
		"nextPage":    nextPage,
		"prevPage":    prevPage,
		"comments":    comments,
	})
}
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
