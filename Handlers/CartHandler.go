package Handlers

import (
	database "ETicaret/Database"
	"ETicaret/Helpers"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

func AddToCart(c *fiber.Ctx) error {
	rdb := database.ConnectRedis()
	productID := c.Params("productID")
	username := Helpers.GetUserName(c)

	existingQuantity, err := rdb.HGet(ctx, username, productID).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		err := c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
		if err != nil {
			return err
		}
	}

	if existingQuantity > 0 {
		err := rdb.HIncrBy(ctx, username, productID, 1).Err()
		if err != nil {
			err := c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
			if err != nil {
				return err
			}
		}
	} else {
		err := rdb.HSet(ctx, username, productID, 1).Err()
		if err != nil {
			err := c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
			if err != nil {
				return err
			}
		}
	}

	return c.JSON(fiber.Map{
		"Message":   "Ürün başarıyla eklendi.",
		"productId": productID,
		"quantity":  1,
	})
}

func RemoveFromCart(c *fiber.Ctx) error {
	productID := c.Params("productID")
	rdb := database.ConnectRedis()
	username := Helpers.GetUserName(c)
	err := rdb.HDel(c.Context(), username, productID).Err()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"Message": "Ürün başarıyla silindi.",
	})
}

func ViewCart(c *fiber.Ctx) error {
	rdb := database.ConnectRedis()
	username := Helpers.GetUserName(c)
	val, err := rdb.HGetAll(context.Background(), username).Result()
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
	username := Helpers.GetUserName(c)

	currentQuantity, err := rdb.HGet(ctx, username, productID).Int()
	if err != nil {
		return err
	}
	if currentQuantity > 1 {
		err := rdb.HIncrBy(ctx, username, productID, -1).Err()
		if err != nil {
			return err
		}
	} else {
		err := rdb.HDel(ctx, username, productID).Err()
		if err != nil {
			return err
		}
	}

	return c.JSON(fiber.Map{
		"Message": "Ürün miktarı başarıyla azaltıldı veya ürün silindi.",
	})
}
