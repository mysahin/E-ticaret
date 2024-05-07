package Controllers

import (
	database "ETicaret/Database"
	"ETicaret/Models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type User struct {
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

func getUserName(c *fiber.Ctx) string {
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
