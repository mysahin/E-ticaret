package Controllers

import (
	database "ETicaret/Database"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

func AddToCart(c *fiber.Ctx) error {
	rdb := database.ConnectRedis()
	productID := c.Params("productID")
	rdb.HIncrBy(ctx, "ProductInCart", productID, 1)

	return c.JSON(fiber.Map{
		"Message":   "Ürün başarıyla eklendi.",
		"productId": productID,
		"quantity":  1,
	})
}

func RemoveFromCart(c *fiber.Ctx) error {
	productID := c.Params("productID")
	rdb := database.ConnectRedis()

	err := rdb.HDel(c.Context(), "ProductInCart", productID).Err()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"Message": "Ürün başarıyla silindi.",
	})
}

func ViewCart(c *fiber.Ctx) error {
	// Redis'ten tüm ürünleri al
	rdb := database.ConnectRedis()

	val, err := rdb.HGetAll(context.Background(), "ProductInCart").Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(val)
	return c.JSON(fiber.Map{
		"Sepetiniz": val,
	})
}

func DecreaseQuantityInCart(c *fiber.Ctx) error {
	productID := c.Params("productID")
	rdb := database.ConnectRedis()

	currentQuantity, err := rdb.HGet(ctx, "ProductInCart", productID).Int()
	if err != nil {
		return err
	}
	if currentQuantity > 1 {
		err := rdb.HIncrBy(ctx, "ProductInCart", productID, -1).Err()
		if err != nil {
			return err
		}
	} else {
		err := rdb.HDel(ctx, "ProductInCart", productID).Err()
		if err != nil {
			return err
		}
	}

	return c.JSON(fiber.Map{
		"Message": "Ürün miktarı başarıyla azaltıldı veya ürün silindi.",
	})
}
