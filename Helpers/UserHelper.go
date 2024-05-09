package Helpers

import (
	"ETicaret/Database"
	"ETicaret/Models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var SecretKey = "secret"

func IsLogin(c *fiber.Ctx) bool {
	db := database.DB.Db
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)

		return false
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user Models.Login
	db.First(&user, "id=?", claims.Issuer)

	return true

}

func getID(c *fiber.Ctx) string {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)

		return ""
	}
	claims := token.Claims.(*jwt.StandardClaims)
	return claims.Issuer
}

func GetUserName(c *fiber.Ctx) string {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)

		return ""
	}
	claims := token.Claims.(*jwt.StandardClaims)
	db := database.DB.Db
	user := new(Models.Login)
	if err := db.First(&user, "id=?", claims.Issuer).Error; err != nil {
		return err.Error()
	}
	return user.UserName
}
