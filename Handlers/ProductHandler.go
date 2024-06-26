package Handlers

import (
	"ETicaret/Controllers"
	database "ETicaret/Database"
	"ETicaret/Helpers"
	"ETicaret/Models"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"strconv"
)

func AddProduct(c *fiber.Ctx) error {
	// Kullanıcı giriş kontrolü
	isLogin := Helpers.IsLogin(c)
	if !isLogin {
		return c.JSON(fiber.Map{
			"Warning": "Önce giriş yapınız!",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Form verileri alınamadı: " + err.Error(),
		})
	}

	addedProduct := new(Models.Product)
	addedProduct.TypeId, err = strconv.ParseUint(form.Value["type_id"][0], 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ürün tipi alınamadı: " + err.Error(),
		})
	}
	addedProduct.ProductName = form.Value["product_name"][0]
	addedProduct.ProductPrice, err = strconv.ParseInt(form.Value["product_price"][0], 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ürün fiyatı alınamadı: " + err.Error(),
		})
	}
	addedProduct.ProductStatement = form.Value["product_statement"][0]
	addedProduct.ProductTitle = form.Value["product_title"][0]
	addedProduct.ProductCount, err = strconv.Atoi(form.Value["product_count"][0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ürün sayısı alınamadı: " + err.Error(),
		})
	}

	// Giriş yapan kullanıcının kullanıcı adını al
	userName := Helpers.GetUserName(c)

	// Yeni ürün oluştur ve veritabanına eklemeden önce ID'yi al
	newProduct := Models.Product{
		TypeId:           addedProduct.TypeId,
		ProductName:      addedProduct.ProductName,
		ProductPrice:     addedProduct.ProductPrice,
		ProductStatement: addedProduct.ProductStatement,
		ProductTitle:     addedProduct.ProductTitle,
		SellerUserName:   userName,
		ProductCount:     addedProduct.ProductCount,
	}
	db := database.DB.Db
	if err := db.Create(&newProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ürün kaydedilirken bir hata oluştu: " + err.Error(),
		})
	}
	if err := db.Find(&newProduct).Error; err != nil {
		return err
	}
	id := int(newProduct.ID)
	fileUrl, errr := Controllers.NewFileController(Uploader, Downloader, BucketName).UploadFile(c, strconv.Itoa(id))
	if errr != nil {
		return c.JSON(fiber.Map{
			"error": errr,
		})
	}
	if err := db.Model(&newProduct).Update("image_url", fileUrl).Error; err != nil {
		return err
	}
	if err := db.Find(&newProduct).Error; err != nil {
		return err
	}
	// Başarılı yanıt gönder
	return c.JSON(fiber.Map{
		"message": "Yeni ürün başarıyla eklendi.",
		"product": newProduct,
	})
}

func ViewMyProduct(c *fiber.Ctx) error {
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

func ViewProductBySeller(c *fiber.Ctx) error {
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
	productSeller := c.Params("seller")
	if err := db.Limit(pageSize).Offset(offset).Where("seller_username = ? AND archived = ?", productSeller, "0").Find(&products).Error; err != nil {
		return err
	}

	var totalRecords int64
	if err := db.Model(&Models.Product{}).Where("seller_username = ? AND archived = ?", productSeller, "0").Count(&totalRecords).Error; err != nil {
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
		"Seller":      productSeller,
		"totalPages":  totalPages,
		"currentPage": currentPage,
		"nextPage":    nextPage,
		"prevPage":    prevPage,
		"products":    products,
	})
}

func ViewProductById(c *fiber.Ctx) error {
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

func ViewProductsByType(c *fiber.Ctx) error {
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

func ViewProductsByCategory(c *fiber.Ctx) error {
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

func DeleteProduct(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		db := database.DB.Db
		productId := c.Params("id")
		var deleteProduct Models.Product
		var file Models.Files

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
			if err := db.First(&file, "product_id = ?", productId).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "dosya bulunamadı",
				})
			}

			Controllers.NewFileController(Uploader, Downloader, BucketName).DeleteFile(c, file.FileName)
			if err := db.Delete(&file).Error; err != nil {
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

func ArchiveProduct(c *fiber.Ctx) error {
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

func EditProduct(c *fiber.Ctx) error {
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

func RateProduct(c *fiber.Ctx) error {
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

func HomePage(c *fiber.Ctx) error {
	searchTerm := c.Query("search")

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

	searchTerm = "%" + searchTerm + "%"
	var products []Models.Product
	if err := db.Limit(pageSize).Offset(offset).Where("product_title ILIKE ? AND archived = ? ", searchTerm, "0").Find(&products).Error; err != nil {
		return err
	}

	var totalRecords int64
	if err := db.Model(&Models.Product{}).Where("product_title ILIKE ? AND archived = ?", searchTerm, "0").Count(&totalRecords).Error; err != nil {
		return err
	}

	var categories []Models.Category
	var types []Models.Type

	if err := db.Where("name ILIKE ?", searchTerm).Find(&categories).Error; err != nil {
		return err
	}

	if err := db.Where("name ILIKE ?", searchTerm).Find(&types).Error; err != nil {
		return err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	currentPage := page

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
		"categories":  categories,
		"types":       types,
	})
}

func SearchPageCategorie(c *fiber.Ctx) error {
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
	searchTerm := c.Query("search")

	if err := db.Where("category_id=?", productCategory).Find(&types).Error; err != nil {
		return err
	}

	for _, x := range types {

		var products []Models.Product
		if err := db.Where("type_id=? AND archived=? AND product_title ILIKE", x.ID, "0", searchTerm).Find(&products).Error; err != nil {
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

func SearchPageType(c *fiber.Ctx) error {
	searchTerm := c.Query("search")
	types := c.Query("type")

	if len(searchTerm) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Arama kelimesi en az 3 harf içermelidir.",
		})
	}

	searchTerm = "%" + searchTerm + "%"

	db := database.DB.Db
	var products []Models.Product
	if err := db.Where("product_title ILIKE ? AND product_type ILIKE ?", searchTerm, types).Find(&products).Error; err != nil {
		products = nil
	}

	return c.JSON(fiber.Map{
		"products": products,
	})
}

func CommentProduct(c *fiber.Ctx) error {
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

var BucketName string = "social-media-mysahin"
var Uploader *s3manager.Uploader
var Downloader *s3.S3
